package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func LecturerRoutes(app *fiber.App, lecturerService servicepostgre.ILecturerService, db *sql.DB) {
	lecturers := app.Group("/api/v1/lecturers", middlewarepostgre.AuthRequired())

	lecturers.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	lecturers.Get("/:id/advisees", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})
}
