package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func LecturerRoutes(app *fiber.App, postgresDB *sql.DB) {
	lecturers := app.Group("/api/v1/lecturers", middlewarepostgre.AuthRequired())

	lecturers.Get("", middlewarepostgre.PermissionRequired(postgresDB, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.GetAllLecturersService(c, postgresDB)
	})

	lecturers.Get("/:id/advisees", middlewarepostgre.PermissionRequired(postgresDB, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.GetLecturerAdviseesService(c, postgresDB)
	})
}

