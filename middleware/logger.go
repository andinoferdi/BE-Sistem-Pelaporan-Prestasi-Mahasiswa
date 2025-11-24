package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	fmt.Printf("[%s] %s %s\n", start.Format("2006-01-02 15:04:05"), c.Method(), c.Path())
	
	err := c.Next()
	
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

