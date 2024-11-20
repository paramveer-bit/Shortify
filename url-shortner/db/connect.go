package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDb(dbName string, collName string) (*mongo.Collection, error) {
	mongoURL := os.Getenv("URL")
	if mongoURL == "" {
		return nil, fmt.Errorf("missing MongoDB URL in environment variables")
	}

	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %v", err)
	}

	collection := client.Database(dbName).Collection(collName)
	fmt.Println("Connected to MongoDB and collection instance created!")

	return collection, nil
}
