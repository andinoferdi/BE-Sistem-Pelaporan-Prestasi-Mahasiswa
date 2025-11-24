package config

import (
	"log"
	"os"
)

func GetLogger() *log.Logger {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	return log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

