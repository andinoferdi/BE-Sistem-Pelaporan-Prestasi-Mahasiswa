package middleware

// #1 proses: import library yang diperlukan untuk fmt, time, dan fiber
import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// #2 proses: middleware untuk logging request method, path, status code, dan duration
func LoggerMiddleware(c *fiber.Ctx) error {
	// #2a proses: catat waktu mulai request
	start := time.Now()
	fmt.Printf("[%s] %s %s\n", start.Format("2006-01-02 15:04:05"), c.Method(), c.Path())

	// #2b proses: lanjutkan ke handler berikutnya
	err := c.Next()

	// #2c proses: hitung duration dan log response dengan status code
	duration := time.Since(start)
	fmt.Printf("[%s] %s %s - %d - %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		c.Method(),
		c.Path(),
		c.Response().StatusCode(),
		duration,
	)

	return err
}
