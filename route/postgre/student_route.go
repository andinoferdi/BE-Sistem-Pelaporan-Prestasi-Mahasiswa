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

// #2 proses: setup semua route untuk student dengan middleware AuthRequired dan PermissionRequired
func StudentRoutes(app *fiber.App, studentService servicepostgre.IStudentService, achievementService servicepostgre.IAchievementService, db *sql.DB) {
	// #2a proses: buat group route untuk students dengan middleware AuthRequired
	students := app.Group("/api/v1/students", middlewarepostgre.AuthRequired())

	// #3 proses: endpoint GET /api/v1/students untuk ambil semua student dengan permission user:manage
	students.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #3a proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #3b proses: panggil service get all students
		studentsList, err := studentService.GetAllStudents(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #3c proses: build response dengan data students
		response := modelpostgre.GetAllStudentsResponse{
			Status: "success",
			Data:   studentsList,
		}

		return c.JSON(response)
	})

	// #4 proses: endpoint GET /api/v1/students/:id untuk ambil student berdasarkan ID dengan permission user:manage
	students.Get("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #4a proses: ambil student ID dari URL parameter dan validasi
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID student wajib diisi.",
			})
		}

		// #4b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #4c proses: panggil service get student by ID
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

		// #4d proses: build response dengan data student
		response := modelpostgre.GetStudentByIDResponse{
			Status: "success",
			Data:   *student,
		}

		return c.JSON(response)
	})

	// #5 proses: endpoint GET /api/v1/students/:id/achievements untuk ambil achievements student dengan pagination
	students.Get("/:id/achievements", middlewarepostgre.PermissionRequired(db, "achievement:read"), func(c *fiber.Ctx) error {
		// #5a proses: ambil student ID dari URL parameter dan validasi
		studentID := c.Params("id")
		if studentID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID student wajib diisi.",
			})
		}

		// #5b proses: ambil dan validasi query parameter page dan limit
		page := helper.GetQueryInt(c, "page", 1)
		limit := helper.GetQueryInt(c, "limit", 10)
		page, limit = helper.ValidatePagination(page, limit)

		// #5c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #5d proses: panggil service get achievements by student ID
		response, err := achievementService.GetAchievementsByStudentID(ctx, studentID, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #5e proses: return response dengan data achievements dan pagination
		return c.JSON(response)
	})

	// #6 proses: endpoint PUT /api/v1/students/:id/advisor untuk update advisor student dengan permission user:manage
	students.Put("/:id/advisor", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #6a proses: ambil student ID dari URL parameter dan validasi
		studentID := c.Params("id")
		if studentID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID student wajib diisi.",
			})
		}

		// #6b proses: parse request body untuk ambil advisor ID
		req := struct {
			AdvisorID string `json:"advisor_id"`
		}{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #6c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #6d proses: panggil service update student advisor
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

		// #6e proses: return response sukses
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Advisor berhasil diupdate",
		})
	})
}
