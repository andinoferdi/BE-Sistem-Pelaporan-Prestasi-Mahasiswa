package route

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

func LecturerRoutes(app *fiber.App, lecturerService servicepostgre.ILecturerService, studentService servicepostgre.IStudentService, db *sql.DB) {
	lecturers := app.Group("/api/v1/lecturers", middlewarepostgre.AuthRequired())

	lecturers.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
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
	})

	lecturers.Get("/:id/advisees", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
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
	})
}
