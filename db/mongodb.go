package db

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var (
	conn                  *mongo.Client // is concurrent safe: https://github.com/mongodb/mongo-go-driver/blob/master/mongo/client.go#L46
	ConnectionTimeout     = 10 * time.Second
	mu                    sync.RWMutex
	SafeDurationInMinutes int = 13
)

const (
	NotFoundDoc = "mongo: no documents in result"
	InvalidPK   = "invalid pk"
	emptyConn   = "Connection is empty"
)

// CheckMongo checks if Mongo is up and running
func CheckMongo() error {
	_, err := CreateMongo()
	if err != nil {
		return err
	}

	return ping()
}

// CreateMongo cria uma nova instancia de conexão com o mongodb
func CreateMongo() (*mongo.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if conn != nil {
		return conn, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	var err error
	l := log.CreateLog()
	conn, err = mongo.Connect(ctx, getClientOptions())
	if err != nil {
		l.Error(err.Error(), "mongodb.CreateMongo - Error creating mongo connection")
		return conn, err
	}

	return conn, nil
}

func getClientOptions() *options.ClientOptions {
	mongoURL := config.Get().MongoURL
	co := options.Client()
	co.SetRetryWrites(true)
	co.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	co.SetConnectTimeout(5 * time.Second)
	co.SetMaxConnIdleTime(10 * time.Second)
	co.SetMaxPoolSize(512)

	if config.Get().ForceTLS {
		co.SetTLSConfig(&tls.Config{})
	}

	return co.ApplyURI(fmt.Sprintf("mongodb://%s", mongoURL)).SetAuth(mongoCredential())
}

func mongoCredential() options.Credential {
	user := config.Get().MongoUser
	password := config.Get().MongoPassword
	var database string
	if config.Get().MongoAuthSource != "" {
		database = config.Get().MongoAuthSource
	} else {
		database = config.Get().MongoDatabase
	}

	credential := options.Credential{
		Username:   user,
		Password:   password,
		AuthSource: database,
	}

	if config.Get().ForceTLS {
		credential.AuthMechanism = "SCRAM-SHA-1"
	}

	return credential
}

//SaveBoleto salva um boleto no mongoDB
func SaveBoleto(boleto models.BoletoView) error {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	l := log.CreateLog()
	conn, err := CreateMongo()
	if err != nil {
		l.Error(err.Error(), fmt.Sprintf("mongodb.CreateMongo - Error creating mongo connection while saving boleto %v", boleto))
		return err
	}

	collection := conn.Database(config.Get().MongoDatabase).Collection(config.Get().MongoBoletoCollection)
	_, err = collection.InsertOne(ctx, boleto)

	return err
}

//GetBoletoByID busca um boleto pelo ID que vem na URL
//O retorno será um objeto BoletoView, o tempo decorrido da operação (em milisegundos) e algum erro ocorrido durante a operação
func GetBoletoByID(id, pk string) (models.BoletoView, int64, error) {
	start := time.Now()

	result := models.BoletoView{}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	l := log.CreateLog()
	conn, err := CreateMongo()
	if err != nil {
		l.Error(err.Error(), fmt.Sprintf("mongodb.GetBoletoByID - Error creating mongo connection for id %s and pk %s", id, pk))
		return result, time.Since(start).Milliseconds(), err
	}
	collection := conn.Database(config.Get().MongoDatabase).Collection(config.Get().MongoBoletoCollection)

	for i := 0; i <= config.Get().RetryNumberGetBoleto; i++ {

		var filter primitive.M
		if len(id) == 24 {
			d, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return result, time.Since(start).Milliseconds(), fmt.Errorf("Error: %s\n", err)
			}
			filter = bson.M{"_id": d}
		} else {
			filter = bson.M{"id": id}
		}
		err = collection.FindOne(ctx, filter).Decode(&result)

		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			continue
		} else {
			break
		}
	}

	if err != nil {
		return models.BoletoView{}, time.Since(start).Milliseconds(), err
	} else if !hasValidKey(result, pk) {
		return models.BoletoView{}, time.Since(start).Milliseconds(), errors.New(InvalidPK)
	}

	// Changing dates as LocalDateTime, in order to keep the same time.Time attributes the mgo used return
	result.Boleto.Title.ExpireDateTime = util.TimeToLocalTime(result.Boleto.Title.ExpireDateTime)
	result.Boleto.Title.CreateDate = util.TimeToLocalTime(result.Boleto.Title.CreateDate)
	result.CreateDate = util.TimeToLocalTime(result.CreateDate)

	return result, time.Since(start).Milliseconds(), nil
}

//GetUserCredentials Busca as Credenciais dos Usuários
func GetUserCredentials() ([]models.Credentials, error) {
	result := []models.Credentials{}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	l := log.CreateLog()
	conn, err := CreateMongo()
	if err != nil {
		l.Error(err.Error(), "mongodb.GetUserCredentials - Error creating mongo connection")
		return result, err
	}
	collection := conn.Database(config.Get().MongoDatabase).Collection(config.Get().MongoCredentialsCollection)

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetTokenByClientIDAndIssuerBank fetches a token by clientID and issuerBank
func GetTokenByClientIDAndIssuerBank(clientID, issuerBank string) (models.Token, error) {
	result := models.Token{}
	if clientID == "" || issuerBank == "" {
		return result, fmt.Errorf("fields clientID and issuerBank cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := CreateMongo()
	if err != nil {
		return result, err
	}

	filter := bson.D{
		primitive.E{Key: "clientid", Value: clientID},
		primitive.E{Key: "issuerbank", Value: issuerBank},
	}
	opts := options.FindOne().SetSort(bson.D{primitive.E{Key: "createdat", Value: -1}})

	collection := conn.Database(config.Get().MongoDatabase).Collection(config.Get().MongoTokenCollection)
	err = collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return result, err
	}

	if time.Now().After(result.CreatedAt.Add(time.Duration(SafeDurationInMinutes) * time.Minute)) { // safety margin
		result = models.Token{} // a little help for GC
	}

	return result, nil
}

// SaveToken saves an access token at mongoDB
func SaveToken(token models.Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := CreateMongo()
	if err != nil {
		return err
	}

	collection := conn.Database(config.Get().MongoDatabase).Collection(config.Get().MongoTokenCollection)
	_, err = collection.InsertOne(ctx, token)

	return err
}

func hasValidKey(r models.BoletoView, pk string) bool {
	return r.SecretKey == "" || r.PublicKey == pk
}

func ping() error {
	if conn == nil {
		return fmt.Errorf(emptyConn)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	err := conn.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	return nil
}
