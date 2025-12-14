package route_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	utilspostgre "sistem-pelaporan-prestasi-mahasiswa/utils/postgre"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
)

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

func createTestToken(userID, email, roleID string) (string, error) {
	user := createTestUser(userID, email, roleID)
	return utilspostgre.GenerateToken(user)
}

func createRequestWithToken(method, path string, body interface{}, token string) *http.Request {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func assertStatusCode(t *testing.T, resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		t.Errorf("Expected status code %d, got %d", expected, resp.StatusCode)
	}
}

func setupTestApp() *fiber.App {
	return fiber.New()
}

func stringPtr(s string) *string {
	return &s
}

func setupTestDBForRoute(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	return db, mock
}

func getPermissionQuery() string {
	return `(?i)SELECT\s+COUNT\(\*\)\s+>\s+0\s+FROM\s+role_permissions\s+rp\s+INNER\s+JOIN\s+permissions\s+p\s+ON\s+rp\.permission_id\s+=\s+p\.id\s+INNER\s+JOIN\s+users\s+u\s+ON\s+u\.role_id\s+=\s+rp\.role_id\s+WHERE\s+u\.id\s+=\s+\$1\s+AND\s+p\.name\s+=\s+\$2`
}
