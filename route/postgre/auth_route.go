package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, db *sql.DB, instanceID string) {
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		c.Locals("server_instance_id", instanceID)
		return servicepostgre.HealthCheckService(c)
	})

	auth := app.Group("/api/v1/auth")

	auth.Post("/login", func(c *fiber.Ctx) error {
		return servicepostgre.LoginService(c, db)
	})

	auth.Post("/refresh", func(c *fiber.Ctx) error {
		return servicepostgre.RefreshTokenService(c, db)
	})

	protected := auth.Group("", middlewarepostgre.AuthRequired())

	protected.Post("/logout", func(c *fiber.Ctx) error {
		return servicepostgre.LogoutService(c, db)
	})

	protected.Get("/profile", func(c *fiber.Ctx) error {
		return servicepostgre.GetProfileService(c, db)
	})
}

