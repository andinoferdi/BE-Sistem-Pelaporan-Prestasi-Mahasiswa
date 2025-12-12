package route

import (
	"context"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	"sistem-pelaporan-prestasi-mahasiswa/helper"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NotificationRoutes(app *fiber.App, notificationService servicepostgre.INotificationService) {
	notifications := app.Group("/api/v1/notifications", middlewarepostgre.AuthRequired())

	notifications.Get("", func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		page := helper.GetQueryInt(c, "page", 1)
		limit := helper.GetQueryInt(c, "limit", 10)
		page, limit = helper.ValidatePagination(page, limit)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := notificationService.GetNotifications(ctx, userID, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error mengambil notifications: " + err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	notifications.Get("/unread-count", func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := notificationService.GetUnreadCount(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error mengambil unread count: " + err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	notifications.Put("/:id/read", func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		notificationID := c.Params("id")
		if notificationID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID notification wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := notificationService.MarkAsRead(ctx, notificationID, userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal mengambil pengguna",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	notifications.Put("/read-all", func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := notificationService.MarkAllAsRead(ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Error menandai semua notification sebagai read: " + err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})
}
