package main

import (
	"log"
	"sistem-pelaporan-prestasi-mahasiswa/config"
	"sistem-pelaporan-prestasi-mahasiswa/database"
)

func main() {
	config.LoadEnv()

	postgresDB := database.ConnectDB()
	defer postgresDB.Close()

	mongoDB := database.ConnectMongoDB()

	if err := database.RunMigrations(postgresDB, mongoDB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}

