package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api")

	api.Post("/login", func(c *fiber.Ctx) error {
		return servicepostgre.LoginService(c, db)
	})

	protected := api.Group("", middlewarepostgre.AuthRequired())

	protected.Get("/profile", func(c *fiber.Ctx) error {
		return servicepostgre.GetProfileService(c, db)
	})
}

