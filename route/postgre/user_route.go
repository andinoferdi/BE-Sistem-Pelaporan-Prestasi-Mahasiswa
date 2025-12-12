package route

import (
	"context"
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userService servicepostgre.IUserService, db *sql.DB) {
	users := app.Group("/api/v1/users", middlewarepostgre.AuthRequired())

	users.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		users, err := userService.GetAllUsers(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(users)
	})

	users.Get("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := userService.GetUserByID(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal mengambil pengguna",
				"message": err.Error(),
			})
		}

		return c.JSON(user)
	})

	users.Post("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	users.Put("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	users.Delete("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	users.Put("/:id/role", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})
}
