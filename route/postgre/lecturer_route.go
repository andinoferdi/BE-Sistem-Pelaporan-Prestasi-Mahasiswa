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

// GetAllLecturers godoc
// @Summary Get all lecturers
// @Description Mengambil daftar semua lecturer dari database. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} model.GetAllLecturersResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /lecturers [get]
func GetAllLecturers(lecturerService servicepostgre.ILecturerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		lecturersList, err := lecturerService.GetAllLecturers(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		response := model.GetAllLecturersResponse{
			Status: "success",
			Data:   lecturersList,
		}

		return c.JSON(response)
	}
}

// GetLecturerAdvisees godoc
// @Summary Get lecturer advisees
// @Description Mengambil daftar mahasiswa bimbingan dari lecturer tertentu. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Lecturer ID"
// @Success 200 {object} model.GetAllStudentsResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /lecturers/{id}/advisees [get]
func GetLecturerAdvisees(lecturerService servicepostgre.ILecturerService, studentService servicepostgre.IStudentService) fiber.Handler {
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

		studentsList, err := studentService.GetStudentsByAdvisorID(ctx, lecturer.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		response := model.GetAllStudentsResponse{
			Status: "success",
			Data:   studentsList,
		}

		return c.JSON(response)
	}
}

// #2 proses: setup semua route untuk lecturer dengan middleware AuthRequired dan PermissionRequired
func LecturerRoutes(app *fiber.App, lecturerService servicepostgre.ILecturerService, studentService servicepostgre.IStudentService, db *sql.DB) {
	lecturers := app.Group("/api/v1/lecturers", middlewarepostgre.AuthRequired())

	lecturers.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), GetAllLecturers(lecturerService))
	lecturers.Get("/:id/advisees", middlewarepostgre.PermissionRequired(db, "user:manage"), GetLecturerAdvisees(lecturerService, studentService))
}
