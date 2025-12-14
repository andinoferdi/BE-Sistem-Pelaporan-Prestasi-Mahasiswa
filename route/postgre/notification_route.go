package route

// #1 proses: import library yang diperlukan untuk context, service, helper, middleware, time, dan fiber
import (
	"context"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	"sistem-pelaporan-prestasi-mahasiswa/helper"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"time"

	"github.com/gofiber/fiber/v2"
)

// #2 proses: setup semua route untuk notifikasi dengan middleware AuthRequired
func NotificationRoutes(app *fiber.App, notificationService servicepostgre.INotificationService) {
	// #2a proses: buat group route untuk notifications dengan middleware AuthRequired
	notifications := app.Group("/api/v1/notifications", middlewarepostgre.AuthRequired())

	// #3 proses: endpoint GET /api/v1/notifications untuk ambil notifikasi user dengan pagination
	notifications.Get("", func(c *fiber.Ctx) error {
		// #3a proses: ambil user ID dari context yang diset oleh middleware
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		// #3b proses: ambil dan validasi query parameter page dan limit
		page := helper.GetQueryInt(c, "page", 1)
		limit := helper.GetQueryInt(c, "limit", 10)
		page, limit = helper.ValidatePagination(page, limit)

		// #3c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #3d proses: panggil service get notifications
		response, err := notificationService.GetNotifications(ctx, userID, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error mengambil notifications: " + err.Error(),
			})
		}

		// #3e proses: return response dengan data notifications dan pagination
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #4 proses: endpoint GET /api/v1/notifications/unread-count untuk ambil jumlah notifikasi belum dibaca
	notifications.Get("/unread-count", func(c *fiber.Ctx) error {
		// #4a proses: ambil user ID dari context yang diset oleh middleware
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		// #4b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #4c proses: panggil service get unread count
		response, err := notificationService.GetUnreadCount(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error mengambil unread count: " + err.Error(),
			})
		}

		// #4d proses: return response dengan count notifikasi belum dibaca
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #5 proses: endpoint PUT /api/v1/notifications/:id/read untuk tandai notifikasi sebagai sudah dibaca
	notifications.Put("/:id/read", func(c *fiber.Ctx) error {
		// #5a proses: ambil user ID dari context yang diset oleh middleware
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		// #5b proses: ambil notification ID dari URL parameter dan validasi
		notificationID := c.Params("id")
		if notificationID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID notification wajib diisi.",
			})
		}

		// #5c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #5d proses: panggil service mark as read
		response, err := notificationService.MarkAsRead(ctx, notificationID, userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal mengambil pengguna",
				"message": err.Error(),
			})
		}

		// #5e proses: return response dengan notifikasi yang sudah diupdate
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #6 proses: endpoint PUT /api/v1/notifications/read-all untuk tandai semua notifikasi sebagai sudah dibaca
	notifications.Put("/read-all", func(c *fiber.Ctx) error {
		// #6a proses: ambil user ID dari context yang diset oleh middleware
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		// #6b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #6c proses: panggil service mark all as read
		response, err := notificationService.MarkAllAsRead(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error menandai semua notification sebagai read: " + err.Error(),
			})
		}

		// #6d proses: return response sukses
		return c.Status(fiber.StatusOK).JSON(response)
	})
}
