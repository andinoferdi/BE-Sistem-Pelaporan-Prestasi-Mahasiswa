// #1 proses: package main sebagai entry point aplikasi
package main

// @title Sistem Pelaporan Prestasi Mahasiswa API
// @version 1.0
// @description API untuk sistem pelaporan prestasi mahasiswa dengan dukungan multi-role (Admin, Dosen Wali, Mahasiswa) dan field prestasi yang fleksibel
// @host localhost:3001
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer {token}"

// #2 proses: import library yang diperlukan untuk log, os, repository, service, config, database, middleware, route, dan uuid
import (
	"log"
	"os"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	"sistem-pelaporan-prestasi-mahasiswa/config"
	configmongo "sistem-pelaporan-prestasi-mahasiswa/config/mongo"
	"sistem-pelaporan-prestasi-mahasiswa/database"
	_ "sistem-pelaporan-prestasi-mahasiswa/docs"
	"sistem-pelaporan-prestasi-mahasiswa/middleware"
	routepostgre "sistem-pelaporan-prestasi-mahasiswa/route/postgre"

	"github.com/google/uuid"
)

// #3 proses: variable global untuk menyimpan server instance ID yang unik
var serverInstanceID string

// #4 proses: fungsi main sebagai entry point aplikasi yang menginisialisasi semua komponen
func main() {
	// #4a proses: load environment variables dari file .env
	config.LoadEnv()

	// #4b proses: generate server instance ID unik menggunakan UUID
	serverInstanceID = uuid.New().String()
	log.Printf("Server instance ID: %s", serverInstanceID)

	// #4c proses: buat directory uploads jika belum ada untuk menyimpan file yang diupload
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	// #4d proses: connect ke PostgreSQL database dan defer close connection
	postgresDB := database.ConnectDB()
	defer postgresDB.Close()

	// #4e proses: connect ke MongoDB database
	mongoDB := database.ConnectMongoDB()

	// #4f proses: inisialisasi Fiber app dengan config MongoDB
	app := configmongo.NewApp()
	app.Use(middleware.LoggerMiddleware)

	// #4f1 proses: setup Swagger UI untuk dokumentasi API dengan config untuk hide models dan tag order
	app.Get("/swagger/*", config.CustomSwaggerHandler())

	// #4g proses: inisialisasi semua repository dengan dependency injection
	userRepo := repositorypostgre.NewUserRepository(postgresDB)
	studentRepo := repositorypostgre.NewStudentRepository(postgresDB)
	lecturerRepo := repositorypostgre.NewLecturerRepository(postgresDB)
	achievementRefRepo := repositorypostgre.NewAchievementReferenceRepository(postgresDB)
	achievementRepo := repositorymongo.NewAchievementRepository(mongoDB)
	notificationRepo := repositorypostgre.NewNotificationRepository(postgresDB)

	// #4h proses: inisialisasi semua service dengan dependency injection dari repository
	authService := servicepostgre.NewAuthService(userRepo)
	userService := servicepostgre.NewUserService(userRepo, studentRepo, lecturerRepo, postgresDB)
	studentService := servicepostgre.NewStudentService(studentRepo, userRepo, lecturerRepo)
	lecturerService := servicepostgre.NewLecturerService(userRepo, lecturerRepo)
	notificationService := servicepostgre.NewNotificationService(notificationRepo, studentRepo, userRepo, achievementRepo)
	achievementService := servicepostgre.NewAchievementService(achievementRepo, achievementRefRepo, userRepo, studentRepo, notificationService)
	reportService := servicepostgre.NewReportService(achievementRepo, achievementRefRepo, studentRepo, userRepo, lecturerRepo)

	// #4i proses: register semua route dengan dependency injection dari service
	routepostgre.AuthRoutes(app, authService, serverInstanceID)
	routepostgre.UserRoutes(app, userService, studentService, lecturerService, postgresDB)
	routepostgre.AchievementRoutes(app, achievementService, postgresDB)
	routepostgre.StudentRoutes(app, studentService, achievementService, postgresDB)
	routepostgre.LecturerRoutes(app, lecturerService, studentService, postgresDB)
	routepostgre.ReportRoutes(app, reportService, postgresDB)
	routepostgre.NotificationRoutes(app, notificationService)

	// #4j proses: ambil port dari environment variable atau gunakan default 3001
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3001"
	}

	// #4k proses: start server pada port yang ditentukan
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

// #5 proses: fungsi helper untuk mendapatkan server instance ID yang sedang berjalan
func GetServerInstanceID() string {
	return serverInstanceID
}
