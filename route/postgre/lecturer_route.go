package route

import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LecturerRoutes(app *fiber.App, lecturerService servicepostgre.ILecturerService, db *sql.DB) {
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
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})
}
