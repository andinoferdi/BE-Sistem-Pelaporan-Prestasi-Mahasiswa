package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() *mongo.Database {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}

	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		databaseName = "sppm_2025"
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	fmt.Println("Successfully connected to MongoDB database:", databaseName)
	return client.Database(databaseName)
}

