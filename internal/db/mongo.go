package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

func DBConnect() error {

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error reading uri for db")
	}
	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		return fmt.Errorf("MONGO_URI is not set")
	}

	client, err := mongo.Connect(
		options.Client().ApplyURI(uri),
	)

	if err != nil {
		return err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return err
	}

	Client = client
	return nil
}
