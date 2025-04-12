package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
	once     sync.Once
)

func Connect() {
	once.Do(func() {
		if os.Getenv("ENV") != "production" {
			err := godotenv.Load(".env")
			if err != nil {
				log.Fatal("Error loading .env file", err)
			}
		}

		MONGODB_URI := os.Getenv("MONGODB_URI")
		clientOptions := options.Client().ApplyURI(MONGODB_URI)

		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatal("Error connecting to DB", err)
		}

		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatal("Error pinging DB", err)
		}

		fmt.Println("Connected to Mongo")

		Client = client
		Database = client.Database("diary")
	})
}

func GetCollection(collectionName string) *mongo.Collection {
	return Database.Collection(collectionName)
}
