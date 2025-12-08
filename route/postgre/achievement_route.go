package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AchievementRoutes(app *fiber.App, postgresDB *sql.DB, mongoDB *mongo.Database) {
	app.Get("/api/v1/achievements/stats", func(c *fiber.Ctx) error {
		return servicepostgre.GetAchievementStatsService(c, postgresDB)
	})

	achievements := app.Group("/api/v1/achievements", middlewarepostgre.AuthRequired())

	achievements.Get("", func(c *fiber.Ctx) error { return servicepostgre.GetAchievementsService(c, postgresDB, mongoDB) })

	achievements.Get("/:id", middlewarepostgre.PermissionRequired(postgresDB, "achievement:read"), func(c *fiber.Ctx) error {
		return servicepostgre.GetAchievementByIDService(c, postgresDB, mongoDB)
	})

	achievements.Post("", middlewarepostgre.PermissionRequired(postgresDB, "achievement:create"), func(c *fiber.Ctx) error {
		return servicepostgre.CreateAchievementService(c, postgresDB, mongoDB)
	})

	achievements.Put("/:id", middlewarepostgre.PermissionRequired(postgresDB, "achievement:update"), func(c *fiber.Ctx) error {
		return servicepostgre.UpdateAchievementService(c, postgresDB, mongoDB)
	})

	achievements.Post("/:id/attachments", middlewarepostgre.PermissionRequired(postgresDB, "achievement:update"), func(c *fiber.Ctx) error {
		return servicepostgre.UploadFileService(c, postgresDB, mongoDB)
	})

	achievements.Post("/:id/submit", middlewarepostgre.PermissionRequired(postgresDB, "achievement:update"), func(c *fiber.Ctx) error {
		return servicepostgre.SubmitAchievementService(c, postgresDB)
	})

	achievements.Post("/:id/verify", middlewarepostgre.PermissionRequired(postgresDB, "achievement:verify"), func(c *fiber.Ctx) error {
		return servicepostgre.VerifyAchievementService(c, postgresDB)
	})

	achievements.Post("/:id/reject", middlewarepostgre.PermissionRequired(postgresDB, "achievement:verify"), func(c *fiber.Ctx) error {
		return servicepostgre.RejectAchievementService(c, postgresDB)
	})

	achievements.Get("/:id/history", middlewarepostgre.PermissionRequired(postgresDB, "achievement:read"), func(c *fiber.Ctx) error {
		return servicepostgre.GetAchievementHistoryService(c, postgresDB)
	})

	achievements.Delete("/:id", middlewarepostgre.PermissionRequired(postgresDB, "achievement:delete"), func(c *fiber.Ctx) error {
		return servicepostgre.DeleteAchievementService(c, postgresDB, mongoDB)
	})
}

