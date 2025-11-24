package config

import (
	"github.com/gofiber/fiber/v2"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Internal server error: " + err.Error(),
			})
		},
	})

	app.Static("/uploads", "./uploads")

	return app
}

