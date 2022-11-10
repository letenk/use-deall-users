package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDB() *mongo.Database {
	dbUser := os.Getenv("MONGO_USER")
	dbPassword := os.Getenv("MONGO_PASSWORD")
	dbHost := os.Getenv("MONGO_HOST")
	dbPort := os.Getenv("MONGO_PORT")
	dbName := os.Getenv("MONGO_DBNAME")

	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPassword, dbHost, dbPort)
	// Create connection options
	clientOptions := options.Client().ApplyURI(dsn)
	// Create new client for connect to mongo
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Panic("MongoDB error connection:", err)
		return nil
	}

	// Connect as run goroutine to monitor status connection
	err = client.Connect(context.TODO())
	if err != nil {
		log.Panic("MongoDB error connection:", err)
		return nil
	}

	log.Println("Connected to mongo!")

	// Return with set database name, and nil for errors
	return client.Database(dbName)
}
