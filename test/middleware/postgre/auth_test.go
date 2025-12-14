package middleware_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	utilspostgre "sistem-pelaporan-prestasi-mahasiswa/utils/postgre"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
)

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	return db, mock
}

func createTestUser(userID, email, roleID string) modelpostgre.User {
	return modelpostgre.User{
		ID:        userID,
		Email:     email,
		RoleID:    roleID,
		Username:  "testuser",
		FullName:  "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestAuthRequired_ValidToken(t *testing.T) {
	app := fiber.New()
	app.Use(middlewarepostgre.AuthRequired())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	user := createTestUser(userID, email, roleID)
	token, err := utilspostgre.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestAuthRequired_MissingToken(t *testing.T) {
	app := fiber.New()
	app.Use(middlewarepostgre.AuthRequired())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthRequired_InvalidTokenFormat(t *testing.T) {
	app := fiber.New()
	app.Use(middlewarepostgre.AuthRequired())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	app := fiber.New()
	app.Use(middlewarepostgre.AuthRequired())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid_token_here")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestRoleRequired_ValidRole(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role_id", "role-id-1")
		return c.Next()
	})
	app.Use(middlewarepostgre.RoleRequired("role-id-1", "role-id-2"))
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

func TestRoleRequired_InvalidRole(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role_id", "role-id-3")
		return c.Next()
	})
	app.Use(middlewarepostgre.RoleRequired("role-id-1", "role-id-2"))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
}

func TestRoleRequired_MissingRoleID(t *testing.T) {
	app := fiber.New()
	app.Use(middlewarepostgre.RoleRequired("role-id-1"))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestPermissionRequired_ValidPermission(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	})
	app.Use(middlewarepostgre.PermissionRequired(db, "achievement:create"))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	userID := "550e8400-e29b-41d4-a716-446655440000"
	permission := "achievement:create"

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	query := `(?i)SELECT\s+COUNT\(\*\)\s+>\s+0\s+FROM\s+role_permissions\s+rp\s+INNER\s+JOIN\s+permissions\s+p\s+ON\s+rp\.permission_id\s+=\s+p\.id\s+INNER\s+JOIN\s+users\s+u\s+ON\s+u\.role_id\s+=\s+rp\.role_id\s+WHERE\s+u\.id\s+=\s+\$1\s+AND\s+p\.name\s+=\s+\$2`
	mock.ExpectQuery(query).
		WithArgs(userID, permission).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPermissionRequired_InvalidPermission(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	})
	app.Use(middlewarepostgre.PermissionRequired(db, "achievement:delete"))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	userID := "550e8400-e29b-41d4-a716-446655440000"
	permission := "achievement:delete"

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(false)

	query := `(?i)SELECT\s+COUNT\(\*\)\s+>\s+0\s+FROM\s+role_permissions\s+rp\s+INNER\s+JOIN\s+permissions\s+p\s+ON\s+rp\.permission_id\s+=\s+p\.id\s+INNER\s+JOIN\s+users\s+u\s+ON\s+u\.role_id\s+=\s+rp\.role_id\s+WHERE\s+u\.id\s+=\s+\$1\s+AND\s+p\.name\s+=\s+\$2`
	mock.ExpectQuery(query).
		WithArgs(userID, permission).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPermissionRequired_MissingUserID(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	app := fiber.New()
	app.Use(middlewarepostgre.PermissionRequired(db, "achievement:create"))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}
