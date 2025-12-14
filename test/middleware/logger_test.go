package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	middleware "sistem-pelaporan-prestasi-mahasiswa/middleware"

	"github.com/gofiber/fiber/v2"
)

func TestLoggerMiddleware_Success(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestLoggerMiddleware_PostRequest(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("Created")
	})

	req := httptest.NewRequest("POST", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestLoggerMiddleware_NotFound(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/notfound", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestLoggerMiddleware_DoesNotBlockRequest(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestLoggerMiddleware_MultipleRequests(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request %d: %v", i+1, err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Request %d: Expected status code %d, got %d", i+1, http.StatusOK, resp.StatusCode)
		}
	}
}
