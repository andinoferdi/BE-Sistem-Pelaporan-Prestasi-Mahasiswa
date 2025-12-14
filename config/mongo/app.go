package config

// #1 proses: import library yang diperlukan untuk fiber dan cors middleware
import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// #2 proses: inisialisasi Fiber app untuk MongoDB dengan CORS, error handler, dan static files
func NewApp() *fiber.App {
	// #2a proses: buat Fiber app dengan body limit dan custom error handler
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Internal server error: " + err.Error(),
			})
		},
	})

	// #2b proses: setup CORS middleware untuk allow request dari frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// #2c proses: setup static file serving untuk uploads directory
	app.Static("/uploads", "./uploads")

	return app
}
