package config

// #1 proses: import library yang diperlukan untuk log dan os
import (
	"log"
	"os"
)

// #2 proses: setup logger dengan file output ke logs/app.log
func GetLogger() *log.Logger {
	// #2a proses: buat directory logs jika belum ada
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// #2b proses: buka atau buat file log dengan mode append
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// #2c proses: buat logger baru dengan prefix dan format timestamp
	return log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}
