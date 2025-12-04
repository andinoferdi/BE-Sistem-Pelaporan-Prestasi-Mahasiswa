package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AchievementRoutes(app *fiber.App, postgresDB *sql.DB, mongoDB *mongo.Database) {
	achievements := app.Group("/api/v1/achievements", middlewarepostgre.AuthRequired())

	achievements.Get("", func(c *fiber.Ctx) error {
		return servicepostgre.GetAchievementsService(c, postgresDB, mongoDB)
	})

	achievements.Post("", middlewarepostgre.PermissionRequired(postgresDB, "achievement:create"), func(c *fiber.Ctx) error {
		return servicepostgre.CreateAchievementService(c, postgresDB, mongoDB)
	})

	achievements.Post("/upload", middlewarepostgre.PermissionRequired(postgresDB, "achievement:create"), func(c *fiber.Ctx) error {
		return servicepostgre.UploadFileService(c)
	})

	achievements.Post("/:id/submit", middlewarepostgre.PermissionRequired(postgresDB, "achievement:update"), func(c *fiber.Ctx) error {
		return servicepostgre.SubmitAchievementService(c, postgresDB)
	})

	achievements.Delete("/:id", middlewarepostgre.PermissionRequired(postgresDB, "achievement:delete"), func(c *fiber.Ctx) error {
		return servicepostgre.DeleteAchievementService(c, postgresDB, mongoDB)
	})
}

