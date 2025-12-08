package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func StudentRoutes(app *fiber.App, postgresDB *sql.DB, mongoDB *mongo.Database) {
	students := app.Group("/api/v1/students", middlewarepostgre.AuthRequired())

	students.Get("", middlewarepostgre.PermissionRequired(postgresDB, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.GetAllStudentsService(c, postgresDB)
	})

	students.Get("/:id", middlewarepostgre.PermissionRequired(postgresDB, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.GetStudentByIDService(c, postgresDB)
	})

	students.Get("/:id/achievements", middlewarepostgre.PermissionRequired(postgresDB, "achievement:read"), func(c *fiber.Ctx) error {
		return servicepostgre.GetStudentAchievementsService(c, postgresDB, mongoDB)
	})

	students.Put("/:id/advisor", middlewarepostgre.PermissionRequired(postgresDB, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.UpdateStudentAdvisorService(c, postgresDB)
	})
}

