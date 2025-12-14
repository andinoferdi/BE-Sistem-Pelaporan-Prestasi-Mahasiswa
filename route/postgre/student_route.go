package route

// #1 proses: import library yang diperlukan untuk context, database, model, service, helper, middleware, strings, time, dan fiber
import (
	"context"
	"database/sql"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	"sistem-pelaporan-prestasi-mahasiswa/helper"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetAllStudents godoc
// @Summary Get all students
// @Description Mengambil daftar semua student dari database. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Students
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} modelpostgre.GetAllStudentsResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /students [get]
func GetAllStudents(studentService servicepostgre.IStudentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		studentsList, err := studentService.GetAllStudents(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		response := modelpostgre.GetAllStudentsResponse{
			Status: "success",
			Data:   studentsList,
		}

		return c.JSON(response)
	}
}

// GetStudentByID godoc
// @Summary Get student by ID
// @Description Mengambil data student spesifik berdasarkan ID. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Students
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Student ID"
// @Success 200 {object} modelpostgre.GetStudentByIDResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /students/{id} [get]
func GetStudentByID(studentService servicepostgre.IStudentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID student wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		student, err := studentService.GetStudentByID(ctx, id)
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

		response := modelpostgre.GetStudentByIDResponse{
			Status: "success",
			Data:   *student,
		}

		return c.JSON(response)
	}
}

// GetStudentAchievements godoc
// @Summary Get student achievements
// @Description Mengambil daftar achievements milik student tertentu dengan pagination. Dapat diakses dengan permission achievement:read
// @Tags Students
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Student ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /students/{id}/achievements [get]
func GetStudentAchievements(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentID := c.Params("id")
		if studentID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID student wajib diisi.",
			})
		}

		page := helper.GetQueryInt(c, "page", 1)
		limit := helper.GetQueryInt(c, "limit", 10)
		page, limit = helper.ValidatePagination(page, limit)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.GetAchievementsByStudentID(ctx, studentID, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// UpdateStudentAdvisor godoc
// @Summary Update student advisor
// @Description Memperbarui dosen wali untuk student tertentu. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Students
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Student ID"
// @Param body body object true "Advisor ID" example({"advisor_id": "uuid-here"})
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /students/{id}/advisor [put]
func UpdateStudentAdvisor(studentService servicepostgre.IStudentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentID := c.Params("id")
		if studentID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID student wajib diisi.",
			})
		}

		req := struct {
			AdvisorID string `json:"advisor_id"`
		}{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := studentService.UpdateStudentAdvisor(ctx, studentID, req.AdvisorID)
		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Data tidak ditemukan",
					"message": err.Error(),
				})
			}
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal mengupdate advisor",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Advisor berhasil diupdate",
		})
	}
}

// #2 proses: setup semua route untuk student dengan middleware AuthRequired dan PermissionRequired
func StudentRoutes(app *fiber.App, studentService servicepostgre.IStudentService, achievementService servicepostgre.IAchievementService, db *sql.DB) {
	students := app.Group("/api/v1/students", middlewarepostgre.AuthRequired())

	students.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), GetAllStudents(studentService))
	students.Get("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), GetStudentByID(studentService))
	students.Get("/:id/achievements", middlewarepostgre.PermissionRequired(db, "achievement:read"), GetStudentAchievements(achievementService))
	students.Put("/:id/advisor", middlewarepostgre.PermissionRequired(db, "user:manage"), UpdateStudentAdvisor(studentService))
}
