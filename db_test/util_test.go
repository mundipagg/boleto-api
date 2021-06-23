package db_test

import (
	"context"
	"fmt"
	"time"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
