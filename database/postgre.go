package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping PostgreSQL database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database")
	return db
}

