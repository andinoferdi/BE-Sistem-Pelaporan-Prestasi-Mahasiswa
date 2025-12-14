package helper

// #1 proses: import library yang diperlukan untuk database, strconv, strings, fiber, dan uuid
import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// #2 proses: build response sukses dengan status dan data
func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

// #3 proses: build response error dengan status dan message
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status": "error",
		"data": fiber.Map{
			"message": message,
		},
	})
}

// #4 proses: build response untuk validation error dengan status 400
func ValidationErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message)
}

// #5 proses: build response untuk unauthorized dengan status 401
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message)
}

// #6 proses: build response untuk not found dengan status 404
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message)
}

// #7 proses: build response untuk internal server error dengan status 500
func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message)
}

// #8 proses: build response untuk forbidden dengan status 403
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message)
}

// #9 proses: build response untuk conflict dengan status 409
func ConflictResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusConflict, message)
}

// #10 proses: build response untuk unprocessable entity dengan status 422
func UnprocessableEntityResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnprocessableEntity, message)
}

// #11 proses: handle database error dan return response yang sesuai
func HandleDatabaseError(c *fiber.Ctx, err error) error {
	// #11a proses: jika error adalah sql.ErrNoRows, return not found response
	if err == sql.ErrNoRows {
		return NotFoundResponse(c, "Data tidak ditemukan di database.")
	}
	// #11b proses: untuk error lain, return internal server error response
	return InternalServerErrorResponse(c, "Error mengakses database. Detail: "+err.Error())
}

// #12 proses: parse string ID menjadi UUID
func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// #13 proses: cek apakah string adalah UUID yang valid
func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// #14 proses: ambil query parameter integer dengan default value
func GetQueryInt(c *fiber.Ctx, key string, defaultValue int) int {
	// #14a proses: ambil value dari query parameter
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	// #14b proses: convert ke integer, return default jika error
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// #15 proses: ambil query parameter string dengan default value
func GetQueryString(c *fiber.Ctx, key string, defaultValue string) string {
	// #15a proses: ambil value dari query parameter, return default jika kosong
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// #16 proses: ambil user ID dari context yang diset oleh middleware
func GetUserIDFromContext(c *fiber.Ctx) (string, bool) {
	userID, ok := c.Locals("user_id").(string)
	return userID, ok
}

// #17 proses: ambil email dari context yang diset oleh middleware
func GetEmailFromContext(c *fiber.Ctx) (string, bool) {
	email, ok := c.Locals("email").(string)
	return email, ok
}

// #18 proses: ambil role ID dari context yang diset oleh middleware
func GetRoleIDFromContext(c *fiber.Ctx) (string, bool) {
	roleID, ok := c.Locals("role_id").(string)
	return roleID, ok
}

// #19 proses: validasi dan set default untuk pagination page dan limit
func ValidatePagination(page, limit int) (int, int) {
	// #19a proses: set page minimum 1
	if page < 1 {
		page = 1
	}
	// #19b proses: set limit minimum 10 dan maksimum 100
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

// #20 proses: hitung offset untuk pagination berdasarkan page dan limit
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// #21 proses: sanitize search string dengan hapus karakter khusus SQL
func SanitizeSearch(search string) string {
	// #21a proses: trim whitespace dan hapus karakter % dan _ untuk prevent SQL injection
	search = strings.TrimSpace(search)
	search = strings.ReplaceAll(search, "%", "")
	search = strings.ReplaceAll(search, "_", "")
	return search
}

// #22 proses: cek apakah string kosong setelah di-trim
func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
