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

var globalInstanceID string

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Mengecek status aplikasi dan mendapatkan instance ID server
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "instanceId"
// @Router /health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"instanceId": globalInstanceID,
	})
}

// Login godoc
// @Summary Login user
// @Description Login menggunakan username/email dan password untuk mendapatkan JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/login [post]
func Login(authService servicepostgre.IAuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token menggunakan refresh token yang valid
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body model.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} model.RefreshTokenResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/refresh [post]
func RefreshToken(authService servicepostgre.IAuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// Logout godoc
// @Summary Logout user
// @Description Logout user dari sistem
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]string "message"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/logout [post]
func Logout(authService servicepostgre.IAuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Mengambil profil user yang sedang login beserta role dan permissions
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} model.GetProfileResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /auth/profile [get]
func GetProfile(authService servicepostgre.IAuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// #2 proses: setup semua route untuk autentikasi
func AuthRoutes(app *fiber.App, authService servicepostgre.IAuthService, instanceID string) {
	globalInstanceID = instanceID

	app.Get("/api/v1/health", HealthCheck)

	auth := app.Group("/api/v1/auth")

	auth.Post("/login", Login(authService))

	auth.Post("/refresh", RefreshToken(authService))

	protected := auth.Group("", middlewarepostgre.AuthRequired())

	protected.Post("/logout", Logout(authService))

	protected.Get("/profile", GetProfile(authService))
}
