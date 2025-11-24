package config

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func NewApp(db *sql.DB) *fiber.App {
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

