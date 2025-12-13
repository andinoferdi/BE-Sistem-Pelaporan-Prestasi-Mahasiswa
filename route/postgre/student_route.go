package route

import (
	"context"
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func StudentRoutes(app *fiber.App, studentService servicepostgre.IStudentService, achievementService servicepostgre.IAchievementService, db *sql.DB) {
	students := app.Group("/api/v1/students", middlewarepostgre.AuthRequired())

	students.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	students.Get("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	students.Get("/:id/achievements", middlewarepostgre.PermissionRequired(db, "achievement:read"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	students.Put("/:id/advisor", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
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
	})
}
