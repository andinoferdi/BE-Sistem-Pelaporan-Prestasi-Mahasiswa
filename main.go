package main

import (
	"log"
	"os"
	"sistem-pelaporan-prestasi-mahasiswa/config"
	configmongo "sistem-pelaporan-prestasi-mahasiswa/config/mongo"
	"sistem-pelaporan-prestasi-mahasiswa/database"
	"sistem-pelaporan-prestasi-mahasiswa/middleware"
	routepostgre "sistem-pelaporan-prestasi-mahasiswa/route/postgre"

	"github.com/google/uuid"
)

var serverInstanceID string

func main() {
	config.LoadEnv()

	serverInstanceID = uuid.New().String()
	log.Printf("Server instance ID: %s", serverInstanceID)

	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	postgresDB := database.ConnectDB()
	defer postgresDB.Close()

	mongoDB := database.ConnectMongoDB()

	app := configmongo.NewApp()
	app.Use(middleware.LoggerMiddleware)

	routepostgre.UserRoutes(app, postgresDB, serverInstanceID)
	routepostgre.AchievementRoutes(app, postgresDB, mongoDB)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3001"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func GetServerInstanceID() string {
	return serverInstanceID
}
