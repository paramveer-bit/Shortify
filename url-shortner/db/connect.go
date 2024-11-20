package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const connectionString = "mongodb+srv://coderbuddy01:Pg100904@cluster0.dpplfgk.mongodb.net/"

func ConnectDb(dbName string, collName string) (*mongo.Collection, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	mongoURL := os.Getenv("URL")
	// mongoURL := connectionString
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
