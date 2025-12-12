package route

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

func AchievementRoutes(app *fiber.App, achievementService servicepostgre.IAchievementService, db *sql.DB) {
	app.Get("/api/v1/achievements/stats", func(c *fiber.Ctx) error {
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
	})

	achievements := app.Group("/api/v1/achievements", middlewarepostgre.AuthRequired())

	achievements.Get("", func(c *fiber.Ctx) error {
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

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := achievementService.GetAchievements(ctx, userID, roleID, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		return c.JSON(response)
	})

	achievements.Get("/:id", middlewarepostgre.PermissionRequired(db, "achievement:read"), func(c *fiber.Ctx) error {
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
	})

	achievements.Post("", middlewarepostgre.PermissionRequired(db, "achievement:create"), func(c *fiber.Ctx) error {
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
	})

	achievements.Put("/:id", middlewarepostgre.PermissionRequired(db, "achievement:update"), func(c *fiber.Ctx) error {
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
	})

	achievements.Post("/:id/attachments", middlewarepostgre.PermissionRequired(db, "achievement:update"), func(c *fiber.Ctx) error {
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
			"application/pdf":       true,
			"image/jpeg":            true,
			"image/png":             true,
			"application/msword":    true,
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
	})

	achievements.Post("/:id/submit", middlewarepostgre.PermissionRequired(db, "achievement:update"), func(c *fiber.Ctx) error {
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
	})

	achievements.Post("/:id/verify", middlewarepostgre.PermissionRequired(db, "achievement:verify"), func(c *fiber.Ctx) error {
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
	})

	achievements.Post("/:id/reject", middlewarepostgre.PermissionRequired(db, "achievement:verify"), func(c *fiber.Ctx) error {
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
	})

	achievements.Get("/:id/history", middlewarepostgre.PermissionRequired(db, "achievement:read"), func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{
			"error":   "Fitur belum diimplementasikan",
			"message": "Fitur ini belum diimplementasikan.",
		})
	})

	achievements.Delete("/:id", middlewarepostgre.PermissionRequired(db, "achievement:delete"), func(c *fiber.Ctx) error {
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
	})
}
