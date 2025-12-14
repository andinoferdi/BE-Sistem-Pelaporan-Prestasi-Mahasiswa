package database

// #1 proses: import library yang diperlukan untuk context, fmt, log, os, time, dan MongoDB driver
import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// #2 proses: buat koneksi ke MongoDB database dengan timeout
func ConnectMongoDB() *mongo.Database {
	// #2a proses: ambil MongoDB URI dari environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}

	// #2b proses: ambil database name dari environment variable atau gunakan default
	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		databaseName = "sppm_2025"
	}

	// #2c proses: setup client options dengan URI
	clientOptions := options.Client().ApplyURI(mongoURI)

	// #2d proses: buat context dengan timeout 10 detik untuk koneksi
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// #2e proses: connect ke MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// #2f proses: ping database untuk verifikasi koneksi berhasil
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	fmt.Println("Successfully connected to MongoDB database:", databaseName)
	return client.Database(databaseName)
}
