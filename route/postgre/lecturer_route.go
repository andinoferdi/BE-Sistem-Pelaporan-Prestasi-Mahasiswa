package route

// #1 proses: import library yang diperlukan untuk context, database, model, service, middleware, strings, time, dan fiber
import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// #2 proses: setup semua route untuk lecturer dengan middleware AuthRequired dan PermissionRequired
func LecturerRoutes(app *fiber.App, lecturerService servicepostgre.ILecturerService, studentService servicepostgre.IStudentService, db *sql.DB) {
	// #2a proses: buat group route untuk lecturers dengan middleware AuthRequired
	lecturers := app.Group("/api/v1/lecturers", middlewarepostgre.AuthRequired())

	// #3 proses: endpoint GET /api/v1/lecturers untuk ambil semua lecturer dengan permission user:manage
	lecturers.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #3a proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #3b proses: panggil service get all lecturers
		lecturersList, err := lecturerService.GetAllLecturers(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #3c proses: build response dengan data lecturers
		response := model.GetAllLecturersResponse{
			Status: "success",
			Data:   lecturersList,
		}

		return c.JSON(response)
	})

	// #4 proses: endpoint GET /api/v1/lecturers/:id/advisees untuk ambil mahasiswa bimbingan lecturer tertentu
	lecturers.Get("/:id/advisees", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #4a proses: ambil lecturer ID dari URL parameter dan validasi
		lecturerID := c.Params("id")
		if lecturerID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID lecturer wajib diisi.",
			})
		}

		// #4b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #4c proses: validasi lecturer ID ada di database
		lecturer, err := lecturerService.GetLecturerByID(ctx, lecturerID)
		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Data tidak ditemukan",
					"message": err.Error(),
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #4d proses: ambil mahasiswa bimbingan dari lecturer
		studentsList, err := studentService.GetStudentsByAdvisorID(ctx, lecturer.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #4e proses: build response dengan data mahasiswa bimbingan
		response := model.GetAllStudentsResponse{
			Status: "success",
			Data:   studentsList,
		}

		return c.JSON(response)
	})
}
