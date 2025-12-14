package route

// #1 proses: import library yang diperlukan untuk context, model, service, middleware, time, dan fiber
import (
	"context"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"time"

	"github.com/gofiber/fiber/v2"
)

// #2 proses: setup semua route untuk autentikasi
func AuthRoutes(app *fiber.App, authService servicepostgre.IAuthService, instanceID string) {
	// #2a proses: endpoint health check untuk cek status aplikasi
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"instanceId": instanceID,
		})
	})

	// #2b proses: buat group route untuk auth endpoints
	auth := app.Group("/api/v1/auth")

	// #3 proses: endpoint POST /api/v1/auth/login untuk login user
	auth.Post("/login", func(c *fiber.Ctx) error {
		// #3a proses: parse request body ke LoginRequest
		loginReq := new(model.LoginRequest)
		if err := c.BodyParser(loginReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #3b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #3c proses: panggil service login
		response, err := authService.Login(ctx, *loginReq)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": err.Error(),
			})
		}

		// #3d proses: return response sukses dengan token dan user data
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #4 proses: endpoint POST /api/v1/auth/refresh untuk refresh access token
	auth.Post("/refresh", func(c *fiber.Ctx) error {
		// #4a proses: parse request body ke RefreshTokenRequest
		req := new(model.RefreshTokenRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #4b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #4c proses: panggil service refresh token
		response, err := authService.RefreshToken(ctx, req.RefreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": err.Error(),
			})
		}

		// #4d proses: return response sukses dengan token baru
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #5 proses: buat protected group dengan middleware AuthRequired untuk endpoint yang perlu autentikasi
	protected := auth.Group("", middlewarepostgre.AuthRequired())

	// #6 proses: endpoint POST /api/v1/auth/logout untuk logout user
	protected.Post("/logout", func(c *fiber.Ctx) error {
		// #6a proses: ambil user ID dari context yang diset oleh middleware
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		// #6b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #6c proses: panggil service logout
		err := authService.Logout(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error menghapus refresh token: " + err.Error(),
			})
		}

		// #6d proses: return response sukses logout
		return c.JSON(fiber.Map{
			"message": "Logout berhasil",
		})
	})

	// #7 proses: endpoint GET /api/v1/auth/profile untuk ambil profil user yang login
	protected.Get("/profile", func(c *fiber.Ctx) error {
		// #7a proses: ambil user ID dari context yang diset oleh middleware
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		// #7b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #7c proses: panggil service get profile
		response, err := authService.GetProfile(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal mengambil pengguna",
				"message": err.Error(),
			})
		}

		// #7d proses: return response dengan data profil user
		return c.Status(fiber.StatusOK).JSON(response)
	})
}
