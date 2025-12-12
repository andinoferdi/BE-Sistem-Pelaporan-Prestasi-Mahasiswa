package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func StudentRoutes(app *fiber.App, studentService servicepostgre.IStudentService, achievementService servicepostgre.IAchievementService, db *sql.DB) {
	students := app.Group("/api/v1/students", middlewarepostgre.AuthRequired())

	students.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	students.Get("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	students.Get("/:id/achievements", middlewarepostgre.PermissionRequired(db, "achievement:read"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	students.Put("/:id/advisor", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})
}
