package route

import (
	"context"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authService servicepostgre.IAuthService, instanceID string) {
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"instanceId": instanceID,
		})
	})

	auth := app.Group("/api/v1/auth")

	auth.Post("/login", func(c *fiber.Ctx) error {
		loginReq := new(model.LoginRequest)
		if err := c.BodyParser(loginReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := authService.Login(ctx, *loginReq)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response) 
	})

	auth.Post("/refresh", func(c *fiber.Ctx) error {
		req := new(model.RefreshTokenRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := authService.RefreshToken(ctx, req.RefreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	protected := auth.Group("", middlewarepostgre.AuthRequired())

	protected.Post("/logout", func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := authService.Logout(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error menghapus refresh token: " + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Logout berhasil",
		})
	})

	protected.Get("/profile", func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := authService.GetProfile(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal mengambil pengguna",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})
}
