package main

import (
	"log"
	"os"
	"sistem-pelaporan-prestasi-mahasiswa/config"
	configmongo "sistem-pelaporan-prestasi-mahasiswa/config/mongo"
	"sistem-pelaporan-prestasi-mahasiswa/database"
	"sistem-pelaporan-prestasi-mahasiswa/middleware"
	routemongo "sistem-pelaporan-prestasi-mahasiswa/route/mongo"
	routepostgre "sistem-pelaporan-prestasi-mahasiswa/route/postgre"
)

func main() {
	config.LoadEnv()
	
	postgresDB := database.ConnectDB()
	defer postgresDB.Close()
	
	mongoDB := database.ConnectMongoDB()
	
	if err := database.RunMigrations(mongoDB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	
	app := configmongo.NewApp()
	app.Use(middleware.LoggerMiddleware)
	
	routepostgre.UserRoutes(app, postgresDB)
	
	routepostgre.AlumniRoutes(app, postgresDB)
	routepostgre.PekerjaanRoutes(app, postgresDB)
	
	routemongo.AlumniRoutes(app)
	routemongo.PekerjaanRoutes(app)
	routemongo.FileRoutes(app)
	
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

