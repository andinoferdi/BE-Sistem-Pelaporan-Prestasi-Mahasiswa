package route

// #1 proses: import library yang diperlukan untuk context, database, fmt, mime, os, path, model, service, helper, middleware, strings, time, dan fiber
import (
	"context"
	"database/sql"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	modelmongo "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	"sistem-pelaporan-prestasi-mahasiswa/helper"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetAchievementStats godoc
// @Summary Get achievement statistics
// @Description Mengambil statistik achievement secara umum (public endpoint, tidak perlu autentikasi)
// @Tags Achievements
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /achievements/stats [get]
func GetAchievementStats(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.GetAchievementStats(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// GetAchievements godoc
// @Summary Get all achievements
// @Description Mengambil daftar achievements dengan pagination dan filtering berdasarkan role user (Mahasiswa: milik sendiri, Dosen Wali: milik mahasiswa bimbingan, Admin: semua)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param status query string false "Filter by status (draft, submitted, verified, rejected)"
// @Param achievementType query string false "Filter by achievement type"
// @Param sortBy query string false "Sort by field (created_at, updated_at, submitted_at, status)"
// @Param sortOrder query string false "Sort order (ASC, DESC)" default(DESC)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /achievements [get]
func GetAchievements(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		page := helper.GetQueryInt(c, "page", 1)
		limit := helper.GetQueryInt(c, "limit", 10)
		page, limit = helper.ValidatePagination(page, limit)

		statusFilter := helper.GetQueryString(c, "status", "")
		achievementTypeFilter := helper.GetQueryString(c, "achievementType", "")
		sortBy := helper.GetQueryString(c, "sortBy", "")
		sortOrder := helper.GetQueryString(c, "sortOrder", "")

		if sortOrder != "" {
			sortOrder = strings.ToUpper(sortOrder)
			if sortOrder != "ASC" && sortOrder != "DESC" {
				sortOrder = "DESC"
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.GetAchievements(ctx, userID, roleID, page, limit, statusFilter, achievementTypeFilter, sortBy, sortOrder)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// GetAchievementByID godoc
// @Summary Get achievement by ID
// @Description Mengambil detail achievement berdasarkan ID. Akses dibatasi berdasarkan role (Mahasiswa: hanya milik sendiri, Dosen Wali: hanya mahasiswa bimbingan, Admin: semua)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Achievement ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /achievements/{id} [get]
func GetAchievementByID(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.GetAchievementByID(ctx, userID, roleID, mongoID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal mengambil pengguna",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// CreateAchievement godoc
// @Summary Create achievement
// @Description Membuat achievement baru. Hanya dapat diakses oleh Mahasiswa dengan permission achievement:create
// @Tags Achievements
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body modelmongo.CreateAchievementRequest true "Achievement data"
// @Success 200 {object} modelmongo.CreateAchievementResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /achievements [post]
func CreateAchievement(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		req := new(modelmongo.CreateAchievementRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.CreateAchievement(ctx, userID, roleID, *req)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal membuat prestasi",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

// UpdateAchievement godoc
// @Summary Update achievement
// @Description Memperbarui achievement. Hanya dapat diakses oleh Mahasiswa pemilik dengan permission achievement:update. Hanya dapat diupdate jika status adalah draft
// @Tags Achievements
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Achievement ID (MongoDB ObjectID)"
// @Param body body modelmongo.UpdateAchievementRequest true "Achievement data to update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /achievements/{id} [put]
func UpdateAchievement(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		req := new(modelmongo.UpdateAchievementRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.UpdateAchievement(ctx, userID, roleID, mongoID, *req)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal mengupdate prestasi",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// UploadAttachment godoc
// @Summary Upload attachment
// @Description Mengupload file attachment untuk achievement. Hanya dapat diakses oleh Mahasiswa pemilik dengan permission achievement:update. Hanya dapat diupload jika status adalah draft. Format file: PDF, JPG, PNG, DOC, DOCX (max 10MB)
// @Tags Achievements
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path string true "Achievement ID (MongoDB ObjectID)"
// @Param file formData file true "Attachment file"
// @Success 200 {object} modelmongo.Attachment
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /achievements/{id}/attachments [post]
func UploadAttachment(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "File wajib diisi.",
			})
		}

		fileType := mime.TypeByExtension(filepath.Ext(file.Filename))
		if fileType == "" {
			fileType = "application/octet-stream"
		}

		allowedTypes := map[string]bool{
			"application/pdf":    true,
			"image/jpeg":         true,
			"image/png":          true,
			"application/msword": true,
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		}

		if !allowedTypes[fileType] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Tipe file tidak diizinkan. Gunakan PDF, JPG, PNG, DOC, atau DOCX.",
			})
		}

		timestamp := time.Now().Unix()
		safeFilename := fmt.Sprintf("%d-%s", timestamp, file.Filename)
		safeFilename = strings.ReplaceAll(safeFilename, " ", "_")

		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error membuat folder uploads: " + err.Error(),
			})
		}

		filePath := filepath.Join(uploadDir, safeFilename)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error menyimpan file: " + err.Error(),
			})
		}

		fileURL := fmt.Sprintf("/uploads/%s", safeFilename)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		attachment, err := achievementService.UploadFile(ctx, userID, roleID, mongoID, file.Filename, fileURL, fileType)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(attachment)
	}
}

// SubmitAchievement godoc
// @Summary Submit achievement for verification
// @Description Submit achievement untuk verifikasi oleh dosen wali. Hanya dapat diakses oleh Mahasiswa pemilik dengan permission achievement:update. Hanya dapat disubmit jika status adalah draft
// @Tags Achievements
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Achievement ID (MongoDB ObjectID)"
// @Success 200 {object} modelpostgre.UpdateAchievementReferenceResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /achievements/{id}/submit [post]
func SubmitAchievement(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.SubmitAchievement(ctx, userID, roleID, mongoID)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal submit prestasi",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

// VerifyAchievement godoc
// @Summary Verify achievement
// @Description Memverifikasi achievement. Hanya dapat diakses oleh Dosen Wali dengan permission achievement:verify. Hanya dapat diverifikasi jika status adalah submitted
// @Tags Achievements
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Achievement ID (MongoDB ObjectID)"
// @Success 200 {object} modelpostgre.VerifyAchievementResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /achievements/{id}/verify [post]
func VerifyAchievement(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.VerifyAchievement(ctx, userID, roleID, mongoID)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal memverifikasi prestasi",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

// RejectAchievement godoc
// @Summary Reject achievement
// @Description Menolak achievement dengan catatan. Hanya dapat diakses oleh Dosen Wali dengan permission achievement:verify. Hanya dapat ditolak jika status adalah submitted
// @Tags Achievements
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Achievement ID (MongoDB ObjectID)"
// @Param body body modelpostgre.RejectAchievementRequest true "Rejection note"
// @Success 200 {object} modelpostgre.RejectAchievementResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /achievements/{id}/reject [post]
func RejectAchievement(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		req := new(modelpostgre.RejectAchievementRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.RejectAchievement(ctx, userID, roleID, mongoID, *req)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal menolak prestasi",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

// GetAchievementHistory godoc
// @Summary Get achievement history
// @Description Mengambil history perubahan status achievement. Dapat diakses dengan permission achievement:read
// @Tags Achievements
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Achievement ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /achievements/{id}/history [get]
func GetAchievementHistory(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.GetAchievementHistory(ctx, userID, roleID, mongoID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal mengambil history",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// DeleteAchievement godoc
// @Summary Delete achievement
// @Description Menghapus achievement (soft delete). Hanya dapat diakses oleh Mahasiswa pemilik dengan permission achievement:delete. Hanya dapat dihapus jika status adalah draft
// @Tags Achievements
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Achievement ID (MongoDB ObjectID)"
// @Success 200 {object} modelmongo.DeleteAchievementResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /achievements/{id} [delete]
func DeleteAchievement(achievementService servicepostgre.IAchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			})
		}

		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.DeleteAchievement(ctx, userID, roleID, mongoID)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal menghapus prestasi",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

// #2 proses: setup semua route untuk achievement dengan middleware AuthRequired, PermissionRequired, dan RoleRequired
func AchievementRoutes(app *fiber.App, achievementService servicepostgre.IAchievementService, db *sql.DB) {
	app.Get("/api/v1/achievements/stats", GetAchievementStats(achievementService))

	achievements := app.Group("/api/v1/achievements", middlewarepostgre.AuthRequired())

	achievements.Get("", GetAchievements(achievementService))

	achievements.Get("/:id", middlewarepostgre.PermissionRequired(db, "achievement:read"), GetAchievementByID(achievementService))

	achievements.Post("", middlewarepostgre.PermissionRequired(db, "achievement:create"), CreateAchievement(achievementService))
	achievements.Put("/:id", middlewarepostgre.PermissionRequired(db, "achievement:update"), UpdateAchievement(achievementService))
	achievements.Post("/:id/attachments", middlewarepostgre.PermissionRequired(db, "achievement:update"), UploadAttachment(achievementService))
	achievements.Post("/:id/submit", middlewarepostgre.PermissionRequired(db, "achievement:update"), SubmitAchievement(achievementService))
	achievements.Post("/:id/verify", middlewarepostgre.PermissionRequired(db, "achievement:verify"), VerifyAchievement(achievementService))
	achievements.Post("/:id/reject", middlewarepostgre.PermissionRequired(db, "achievement:verify"), RejectAchievement(achievementService))
	achievements.Get("/:id/history", middlewarepostgre.PermissionRequired(db, "achievement:read"), GetAchievementHistory(achievementService))
	achievements.Delete("/:id", middlewarepostgre.PermissionRequired(db, "achievement:delete"), DeleteAchievement(achievementService))
}
