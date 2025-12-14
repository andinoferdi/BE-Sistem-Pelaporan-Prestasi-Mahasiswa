package config

// #1 proses: import library yang diperlukan untuk log dan godotenv
import (
	"log"

	"github.com/joho/godotenv"
)

// #2 proses: load environment variables dari file .env
func LoadEnv() {
	// #2a proses: load file .env menggunakan godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
