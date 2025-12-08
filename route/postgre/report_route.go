package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(app *fiber.App, postgresDB *sql.DB) {
	reports := app.Group("/api/v1/reports", middlewarepostgre.AuthRequired())

	reports.Get("/statistics", func(c *fiber.Ctx) error {
		return servicepostgre.GetStatisticsService(c, postgresDB)
	})

	reports.Get("/student/:id", middlewarepostgre.PermissionRequired(postgresDB, "achievement:read"), func(c *fiber.Ctx) error {
		return servicepostgre.GetStudentReportService(c, postgresDB)
	})
}

