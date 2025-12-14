package database

// #1 proses: import library yang diperlukan untuk database, fmt, log, os, dan driver PostgreSQL
import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// #2 proses: buat koneksi ke PostgreSQL database
func ConnectDB() *sql.DB {
	// #2a proses: ambil DSN dari environment variable
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	// #2b proses: buka koneksi ke database PostgreSQL
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL database:", err)
	}

	// #2c proses: ping database untuk verifikasi koneksi berhasil
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping PostgreSQL database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database")
	return db
}
