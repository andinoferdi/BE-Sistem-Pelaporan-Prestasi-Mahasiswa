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

// #2 proses: setup semua route untuk achievement dengan middleware AuthRequired, PermissionRequired, dan RoleRequired
func AchievementRoutes(app *fiber.App, achievementService servicepostgre.IAchievementService, db *sql.DB) {
	// #2a proses: endpoint GET /api/v1/achievements/stats untuk ambil statistik achievement, public endpoint
	app.Get("/api/v1/achievements/stats", func(c *fiber.Ctx) error {
		// #2b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #2c proses: panggil service get achievement stats
		response, err := achievementService.GetAchievementStats(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #2d proses: return response dengan data statistik
		return c.JSON(response)
	})

	// #2e proses: buat group route untuk achievements dengan middleware AuthRequired
	achievements := app.Group("/api/v1/achievements", middlewarepostgre.AuthRequired())

	// #3 proses: endpoint GET /api/v1/achievements untuk ambil achievements dengan pagination dan filtering
	achievements.Get("", func(c *fiber.Ctx) error {
		// #3a proses: ambil user ID dan role ID dari context yang diset oleh middleware
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

		// #3b proses: ambil dan validasi query parameter untuk pagination
		page := helper.GetQueryInt(c, "page", 1)
		limit := helper.GetQueryInt(c, "limit", 10)
		page, limit = helper.ValidatePagination(page, limit)

		// #3c proses: ambil query parameter untuk filtering dan sorting
		statusFilter := helper.GetQueryString(c, "status", "")
		achievementTypeFilter := helper.GetQueryString(c, "achievementType", "")
		sortBy := helper.GetQueryString(c, "sortBy", "")
		sortOrder := helper.GetQueryString(c, "sortOrder", "")

		// #3d proses: validasi dan normalize sortOrder
		if sortOrder != "" {
			sortOrder = strings.ToUpper(sortOrder)
			if sortOrder != "ASC" && sortOrder != "DESC" {
				sortOrder = "DESC"
			}
		}

		// #3e proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #3f proses: panggil service get achievements dengan semua filter
		response, err := achievementService.GetAchievements(ctx, userID, roleID, page, limit, statusFilter, achievementTypeFilter, sortBy, sortOrder)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #3g proses: return response dengan data achievements dan pagination
		return c.JSON(response)
	})

	// #4 proses: endpoint GET /api/v1/achievements/:id untuk ambil achievement detail dengan permission achievement:read
	achievements.Get("/:id", middlewarepostgre.PermissionRequired(db, "achievement:read"), func(c *fiber.Ctx) error {
		// #4a proses: ambil user ID dan role ID dari context yang diset oleh middleware
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

		// #4b proses: ambil mongo ID dari URL parameter dan validasi
		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		// #4c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #4d proses: panggil service get achievement by ID
		response, err := achievementService.GetAchievementByID(ctx, userID, roleID, mongoID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal mengambil pengguna",
				"message": err.Error(),
			})
		}

		// #4e proses: return response dengan data achievement detail
		return c.JSON(response)
	})

	// #5 proses: endpoint POST /api/v1/achievements untuk buat achievement baru dengan permission achievement:create
	achievements.Post("", middlewarepostgre.PermissionRequired(db, "achievement:create"), func(c *fiber.Ctx) error {
		// #5a proses: ambil user ID dan role ID dari context yang diset oleh middleware
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

		// #5b proses: parse request body ke CreateAchievementRequest
		req := new(modelmongo.CreateAchievementRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #5c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #5d proses: panggil service create achievement
		response, err := achievementService.CreateAchievement(ctx, userID, roleID, *req)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal membuat prestasi",
				"message": err.Error(),
			})
		}

		// #5e proses: return response dengan data achievement yang baru dibuat
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #6 proses: endpoint PUT /api/v1/achievements/:id untuk update achievement dengan permission achievement:update
	achievements.Put("/:id", middlewarepostgre.PermissionRequired(db, "achievement:update"), func(c *fiber.Ctx) error {
		// #6a proses: ambil user ID dan role ID dari context yang diset oleh middleware
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

		// #6b proses: ambil mongo ID dari URL parameter dan validasi
		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		// #6c proses: parse request body ke UpdateAchievementRequest
		req := new(modelmongo.UpdateAchievementRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #6d proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #6e proses: panggil service update achievement
		response, err := achievementService.UpdateAchievement(ctx, userID, roleID, mongoID, *req)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal mengupdate prestasi",
				"message": err.Error(),
			})
		}

		// #6f proses: return response dengan data achievement yang sudah diupdate
		return c.JSON(response)
	})

	// #7 proses: endpoint POST /api/v1/achievements/:id/attachments untuk upload file attachment dengan permission achievement:update
	achievements.Post("/:id/attachments", middlewarepostgre.PermissionRequired(db, "achievement:update"), func(c *fiber.Ctx) error {
		// #7a proses: ambil user ID dan role ID dari context yang diset oleh middleware
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

		// #7b proses: ambil mongo ID dari URL parameter dan validasi
		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		// #7c proses: ambil file dari form data
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "File wajib diisi.",
			})
		}

		// #7d proses: deteksi file type berdasarkan extension
		fileType := mime.TypeByExtension(filepath.Ext(file.Filename))
		if fileType == "" {
			fileType = "application/octet-stream"
		}

		// #7e proses: validasi file type harus salah satu yang diizinkan
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

		// #7f proses: generate safe filename dengan timestamp untuk prevent conflict
		timestamp := time.Now().Unix()
		safeFilename := fmt.Sprintf("%d-%s", timestamp, file.Filename)
		safeFilename = strings.ReplaceAll(safeFilename, " ", "_")

		// #7g proses: buat directory uploads jika belum ada
		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error membuat folder uploads: " + err.Error(),
			})
		}

		// #7h proses: simpan file ke filesystem
		filePath := filepath.Join(uploadDir, safeFilename)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error menyimpan file: " + err.Error(),
			})
		}

		// #7i proses: buat file URL untuk disimpan ke database
		fileURL := fmt.Sprintf("/uploads/%s", safeFilename)

		// #7j proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #7k proses: panggil service upload file
		attachment, err := achievementService.UploadFile(ctx, userID, roleID, mongoID, file.Filename, fileURL, fileType)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #7l proses: return response dengan data attachment
		return c.JSON(attachment)
	})

	// #8 proses: endpoint POST /api/v1/achievements/:id/submit untuk submit achievement untuk verifikasi dengan permission achievement:update
	achievements.Post("/:id/submit", middlewarepostgre.PermissionRequired(db, "achievement:update"), func(c *fiber.Ctx) error {
		// #8a proses: ambil user ID dan role ID dari context yang diset oleh middleware
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

		// #8b proses: ambil mongo ID dari URL parameter dan validasi
		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		// #8c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #8d proses: panggil service submit achievement
		response, err := achievementService.SubmitAchievement(ctx, userID, roleID, mongoID)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal submit prestasi",
				"message": err.Error(),
			})
		}

		// #8e proses: return response dengan data achievement yang sudah di-submit
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #9 proses: endpoint POST /api/v1/achievements/:id/verify untuk verifikasi achievement dengan permission achievement:verify
	achievements.Post("/:id/verify", middlewarepostgre.PermissionRequired(db, "achievement:verify"), func(c *fiber.Ctx) error {
		// #9a proses: ambil user ID dan role ID dari context yang diset oleh middleware
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

		// #9b proses: ambil mongo ID dari URL parameter dan validasi
		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		// #9c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #9d proses: panggil service verify achievement
		response, err := achievementService.VerifyAchievement(ctx, userID, roleID, mongoID)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal memverifikasi prestasi",
				"message": err.Error(),
			})
		}

		// #9e proses: return response dengan data achievement yang sudah di-verify
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #10 proses: endpoint POST /api/v1/achievements/:id/reject untuk tolak achievement dengan permission achievement:verify
	achievements.Post("/:id/reject", middlewarepostgre.PermissionRequired(db, "achievement:verify"), func(c *fiber.Ctx) error {
		// #10a proses: ambil user ID dan role ID dari context yang diset oleh middleware
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

		// #10b proses: ambil mongo ID dari URL parameter dan validasi
		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		// #10c proses: parse request body ke RejectAchievementRequest untuk ambil rejection note
		req := new(modelpostgre.RejectAchievementRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #10d proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #10e proses: panggil service reject achievement
		response, err := achievementService.RejectAchievement(ctx, userID, roleID, mongoID, *req)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal menolak prestasi",
				"message": err.Error(),
			})
		}

		// #10f proses: return response dengan data achievement yang sudah di-reject
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #11 proses: endpoint GET /api/v1/achievements/:id/history untuk ambil history perubahan status achievement dengan permission achievement:read
	achievements.Get("/:id/history", middlewarepostgre.PermissionRequired(db, "achievement:read"), func(c *fiber.Ctx) error {
		// #11a proses: ambil user ID dan role ID dari context yang diset oleh middleware
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

		// #11b proses: ambil mongo ID dari URL parameter dan validasi
		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		// #11c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #11d proses: panggil service get achievement history
		response, err := achievementService.GetAchievementHistory(ctx, userID, roleID, mongoID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal mengambil history",
				"message": err.Error(),
			})
		}

		// #11e proses: return response dengan data history perubahan status
		return c.JSON(response)
	})

	// #12 proses: endpoint DELETE /api/v1/achievements/:id untuk hapus achievement dengan permission achievement:delete
	achievements.Delete("/:id", middlewarepostgre.PermissionRequired(db, "achievement:delete"), func(c *fiber.Ctx) error {
		// #12a proses: ambil user ID dan role ID dari context yang diset oleh middleware
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

		// #12b proses: ambil mongo ID dari URL parameter dan validasi
		mongoID := c.Params("id")
		if mongoID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID prestasi wajib diisi.",
			})
		}

		// #12c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #12d proses: panggil service delete achievement
		response, err := achievementService.DeleteAchievement(ctx, userID, roleID, mongoID)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal menghapus prestasi",
				"message": err.Error(),
			})
		}

		// #12e proses: return response sukses
		return c.Status(fiber.StatusOK).JSON(response)
	})
}
