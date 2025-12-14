package config

// #1 proses: import library yang diperlukan untuk database dan fiber
import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// #2 proses: inisialisasi Fiber app untuk PostgreSQL dengan custom error handler
func NewApp(db *sql.DB) *fiber.App {
	// #2a proses: buat Fiber app dengan custom error handler untuk handle error global
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Internal server error: " + err.Error(),
			})
		},
	})
	return app
}
