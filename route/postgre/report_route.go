package route

// #1 proses: import library yang diperlukan untuk context, database, service, middleware, time, dan fiber
import (
	"context"
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"time"

	"github.com/gofiber/fiber/v2"
)

// #2 proses: setup semua route untuk laporan dengan middleware AuthRequired
func ReportRoutes(app *fiber.App, reportService servicepostgre.IReportService, db *sql.DB) {
	// #2a proses: buat group route untuk reports dengan middleware AuthRequired
	reports := app.Group("/api/v1/reports", middlewarepostgre.AuthRequired())

	// #3 proses: endpoint GET /api/v1/reports/statistics untuk ambil statistik achievement
	reports.Get("/statistics", func(c *fiber.Ctx) error {
		// #3a proses: ambil user ID dan role ID dari context yang diset oleh middleware
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		// #3b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #3c proses: panggil service get statistics
		response, err := reportService.GetStatistics(ctx, userID, roleID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #3d proses: return response dengan data statistik
		return c.JSON(response)
	})

	// #4 proses: endpoint GET /api/v1/reports/student untuk ambil laporan student yang sedang login
	reports.Get("/student", func(c *fiber.Ctx) error {
		// #4a proses: ambil user ID dari context yang diset oleh middleware
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		// #4b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #4c proses: panggil service get current student report
		response, err := reportService.GetCurrentStudentReport(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #4d proses: return response dengan data laporan student
		return c.JSON(response)
	})

	// #5 proses: endpoint GET /api/v1/reports/student/:id untuk ambil laporan student berdasarkan ID
	reports.Get("/student/:id", func(c *fiber.Ctx) error {
		// #5a proses: ambil student ID dari URL parameter dan validasi
		studentID := c.Params("id")
		if studentID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID student wajib diisi.",
			})
		}

		// #5b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #5c proses: panggil service get student report
		response, err := reportService.GetStudentReport(ctx, studentID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #5d proses: return response dengan data laporan student
		return c.JSON(response)
	})

	// #6 proses: endpoint GET /api/v1/reports/lecturer untuk ambil laporan lecturer yang sedang login
	reports.Get("/lecturer", func(c *fiber.Ctx) error {
		// #6a proses: ambil user ID dari context yang diset oleh middleware
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		// #6b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #6c proses: panggil service get current lecturer report
		response, err := reportService.GetCurrentLecturerReport(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #6d proses: return response dengan data laporan lecturer
		return c.JSON(response)
	})

	// #7 proses: endpoint GET /api/v1/reports/lecturer/:id untuk ambil laporan lecturer berdasarkan ID
	reports.Get("/lecturer/:id", func(c *fiber.Ctx) error {
		// #7a proses: ambil lecturer ID dari URL parameter dan validasi
		lecturerID := c.Params("id")
		if lecturerID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID lecturer wajib diisi.",
			})
		}

		// #7b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #7c proses: panggil service get lecturer report
		response, err := reportService.GetLecturerReport(ctx, lecturerID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #7d proses: return response dengan data laporan lecturer
		return c.JSON(response)
	})
}
