package db

import (
	// this are the inbuilt modules
	"context"
	"log"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var Client = ConnectMongo(context.TODO())

func ConnectMongo(ctx context.Context) *mongo.Client {
	// if there's an environment variable available for the mongo url, use it.
	mongoHost := os.Getenv("MONGO_URL")
	if mongoHost == "" {
		//mongoHost = "mongodb://mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=clipset"
		mongoHost = "mongodb://localhost:27017"

	}

	client, err := mongo.NewClient(mongoHost)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// return the connected client
	return client
}

func Collection(client *mongo.Client, database string, collection string) *mongo.Collection {
	return client.Database(database).Collection(collection)
}
