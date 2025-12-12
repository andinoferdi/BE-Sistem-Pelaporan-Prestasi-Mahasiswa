package route

import (
	"context"
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(app *fiber.App, reportService servicepostgre.IReportService, db *sql.DB) {
	reports := app.Group("/api/v1/reports", middlewarepostgre.AuthRequired())

	reports.Get("/statistics", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := reportService.GetStatistics(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	})

	reports.Get("/student/:id", middlewarepostgre.PermissionRequired(db, "achievement:read"), func(c *fiber.Ctx) error {
		studentID := c.Params("id")
		if studentID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID student wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := reportService.GetStudentReport(ctx, studentID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	})
}
