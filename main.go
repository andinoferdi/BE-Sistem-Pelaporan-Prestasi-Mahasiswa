package main

import (
	"log"
	"os"
	"sistem-pelaporan-prestasi-mahasiswa/config"
	configmongo "sistem-pelaporan-prestasi-mahasiswa/config/mongo"
	"sistem-pelaporan-prestasi-mahasiswa/database"
	"sistem-pelaporan-prestasi-mahasiswa/middleware"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
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

	userRepo := repositorypostgre.NewUserRepository(postgresDB)
	studentRepo := repositorypostgre.NewStudentRepository(postgresDB)
	achievementRefRepo := repositorypostgre.NewAchievementReferenceRepository(postgresDB)
	achievementRepo := repositorymongo.NewAchievementRepository(mongoDB)
	notificationRepo := repositorypostgre.NewNotificationRepository(postgresDB)

	authService := servicepostgre.NewAuthService(userRepo)
	userService := servicepostgre.NewUserService(userRepo)
	studentService := servicepostgre.NewStudentService(studentRepo)
	lecturerService := servicepostgre.NewLecturerService(userRepo)
	notificationService := servicepostgre.NewNotificationService(notificationRepo, studentRepo, userRepo, achievementRepo)
	achievementService := servicepostgre.NewAchievementService(achievementRepo, achievementRefRepo, userRepo, studentRepo, notificationService)
	reportService := servicepostgre.NewReportService()

	routepostgre.AuthRoutes(app, authService, serverInstanceID)
	routepostgre.UserRoutes(app, userService, postgresDB)
	routepostgre.AchievementRoutes(app, achievementService, postgresDB)
	routepostgre.StudentRoutes(app, studentService, achievementService, postgresDB)
	routepostgre.LecturerRoutes(app, lecturerService, postgresDB)
	routepostgre.ReportRoutes(app, reportService, postgresDB)
	routepostgre.NotificationRoutes(app, notificationService)

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
