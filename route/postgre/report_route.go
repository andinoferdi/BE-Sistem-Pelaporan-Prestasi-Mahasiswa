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

// GetStatistics godoc
// @Summary Get achievement statistics
// @Description Mengambil statistik achievement berdasarkan role (Mahasiswa: milik sendiri, Dosen Wali: mahasiswa bimbingan, Admin: semua)
// @Tags Reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /reports/statistics [get]
func GetStatistics(reportService servicepostgre.IReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := reportService.GetStatistics(ctx, userID, roleID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// GetCurrentStudentReport godoc
// @Summary Get current student report
// @Description Mengambil laporan prestasi student yang sedang login
// @Tags Reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /reports/student [get]
func GetCurrentStudentReport(reportService servicepostgre.IReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := reportService.GetCurrentStudentReport(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// GetStudentReport godoc
// @Summary Get student report by ID
// @Description Mengambil laporan prestasi student berdasarkan ID
// @Tags Reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /reports/student/{id} [get]
func GetStudentReport(reportService servicepostgre.IReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// GetCurrentLecturerReport godoc
// @Summary Get current lecturer report
// @Description Mengambil laporan prestasi mahasiswa bimbingan untuk lecturer yang sedang login
// @Tags Reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /reports/lecturer [get]
func GetCurrentLecturerReport(reportService servicepostgre.IReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := reportService.GetCurrentLecturerReport(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// GetLecturerReport godoc
// @Summary Get lecturer report by ID
// @Description Mengambil laporan prestasi mahasiswa bimbingan untuk lecturer berdasarkan ID
// @Tags Reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Lecturer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /reports/lecturer/{id} [get]
func GetLecturerReport(reportService servicepostgre.IReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lecturerID := c.Params("id")
		if lecturerID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID lecturer wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := reportService.GetLecturerReport(ctx, lecturerID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// #2 proses: setup semua route untuk laporan dengan middleware AuthRequired
func ReportRoutes(app *fiber.App, reportService servicepostgre.IReportService, db *sql.DB) {
	reports := app.Group("/api/v1/reports", middlewarepostgre.AuthRequired())

	reports.Get("/statistics", GetStatistics(reportService))
	reports.Get("/student", GetCurrentStudentReport(reportService))
	reports.Get("/student/:id", GetStudentReport(reportService))
	reports.Get("/lecturer", GetCurrentLecturerReport(reportService))
	reports.Get("/lecturer/:id", GetLecturerReport(reportService))
}
