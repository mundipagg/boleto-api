package db_test

import (
	"context"
	"fmt"
	"time"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/test"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type stubBoletoRequestDb struct {
	test.StubBoletoRequest
}

func newStubBoletoRequestDb(bank models.BankNumber) *stubBoletoRequestDb {
	expirationDate := time.Now().Add(5 * 24 * time.Hour)

	base := test.NewStubBoletoRequest(bank)
	s := &stubBoletoRequestDb{
		StubBoletoRequest: *base,
	}

	s.Agreement = models.Agreement{
		AgreementNumber: 123456,
		Agency:          "1234",
	}

	s.Title = models.Title{
		ExpireDateTime: expirationDate,
		ExpireDate:     expirationDate.Format("2006-01-02"),
		OurNumber:      12345678901234,
		AmountInCents:  200,
		DocumentNumber: "1234567890A",
		Instructions:   "Campo de instruções -  max 40 caracteres",
		BoletoType:     "OUT",
		BoletoTypeCode: "99",
	}

	s.Recipient = models.Recipient{
		Document: models.Document{
			Type:   "CNPJ",
			Number: "12123123000112",
		},
	}

	s.Buyer = models.Buyer{
		Name: "Willian Jadson Bezerra Menezes Tupinambá",
		Document: models.Document{
			Type:   "CPF",
			Number: "12312312312",
		},
		Address: models.Address{
			Street:     "Rua da Assunção de Sá",
			Number:     "123",
			Complement: "Seção A, s 02",
			ZipCode:    "20520051",
			City:       "Belém do Pará",
			District:   "Açaí",
			StateCode:  "PA",
		},
	}
	return s
}

// deleteBoletoById deleta um boleto no mongoDB
// Usado apenas para fins de teste
func deleteBoletoById(id string) error {
	if id == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	var filter primitive.M
	if len(id) == 24 {
		d, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("Error converting string to objectID: %s\n", err)
		}
		filter = bson.M{"_id": d}
	} else {
		filter = bson.M{"id": id}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbSession, err := db.CreateMongo()
	if err != nil {
		return err
	}

	d := config.Get().MongoDatabase
	cl := config.Get().MongoBoletoCollection
	collection := dbSession.Database(d).Collection(cl)

	_, err = collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil
}

// saveCredential salva uma credencial no mongoDB
// Usado apenas para fins de teste
func saveCredential(credential models.Credentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbSession, err := db.CreateMongo()
	if err != nil {
		return err
	}

	collection := dbSession.Database(config.Get().MongoDatabase).Collection(config.Get().MongoCredentialsCollection)
	_, err = collection.InsertOne(ctx, credential)

	return err

}

// deleteCredentialById deleta uma credencial no mongoDB
// Usado apenas para fins de teste
func deleteCredentialById(id string) error {
	if id == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	var filter primitive.M
	if len(id) == 24 {
		d, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("Error converting string to objectID: %s\n", err)
		}
		filter = bson.M{"_id": d}
	} else {
		filter = bson.M{"id": id}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbSession, err := db.CreateMongo()
	if err != nil {
		return err
	}

	collection := dbSession.Database(config.Get().MongoDatabase).Collection(config.Get().MongoCredentialsCollection)
	_, err = collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil
}

// getUserCredentialByID busca uma credencial pelo ID
// método apenas para fim de teste
func getUserCredentialByID(id string) (models.Credentials, error) {
	result := models.Credentials{}
	if id == "" {
		return result, fmt.Errorf("ID cannot be empty")
	}

	var filter primitive.M
	if len(id) == 24 {
		d, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return result, fmt.Errorf("Error converting string to objectID: %s\n", err)
		}
		filter = bson.M{"_id": d}
	} else {
		filter = bson.M{"id": id}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbSession, err := db.CreateMongo()
	if err != nil {
		return result, err
	}

	collection := dbSession.Database(config.Get().MongoDatabase).Collection(config.Get().MongoCredentialsCollection)
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func deleteTokenByIssuerBank(issuerBank string) error {
	if issuerBank == "" {
		return fmt.Errorf("issuerBank cannot be empty")
	}

	filter := bson.M{"issuerbank": issuerBank}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbSession, err := db.CreateMongo()
	if err != nil {
		return err
	}

	collection := dbSession.Database(config.Get().MongoDatabase).Collection(config.Get().MongoTokenCollection)
	_, err = collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil
}
