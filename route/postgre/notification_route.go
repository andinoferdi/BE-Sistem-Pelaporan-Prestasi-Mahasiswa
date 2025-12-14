package route

// #1 proses: import library yang diperlukan untuk context, service, helper, middleware, time, dan fiber
import (
	"context"
	_ "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	"sistem-pelaporan-prestasi-mahasiswa/helper"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetNotifications godoc
// @Summary Get notifications
// @Description Mengambil daftar notifikasi user dengan pagination
// @Tags Notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /notifications [get]
func GetNotifications(notificationService servicepostgre.INotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// GetUnreadCount godoc
// @Summary Get unread notification count
// @Description Mengambil jumlah notifikasi yang belum dibaca untuk user yang sedang login
// @Tags Notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /notifications/unread-count [get]
func GetUnreadCount(notificationService servicepostgre.INotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// MarkNotificationAsRead godoc
// @Summary Mark notification as read
// @Description Menandai notifikasi sebagai sudah dibaca berdasarkan ID
// @Tags Notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /notifications/{id}/read [put]
func MarkNotificationAsRead(notificationService servicepostgre.INotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// MarkAllNotificationsAsRead godoc
// @Summary Mark all notifications as read
// @Description Menandai semua notifikasi user sebagai sudah dibaca
// @Tags Notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /notifications/read-all [put]
func MarkAllNotificationsAsRead(notificationService servicepostgre.INotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// #2 proses: setup semua route untuk notifikasi dengan middleware AuthRequired
func NotificationRoutes(app *fiber.App, notificationService servicepostgre.INotificationService) {
	notifications := app.Group("/api/v1/notifications", middlewarepostgre.AuthRequired())

	notifications.Get("", GetNotifications(notificationService))
	notifications.Get("/unread-count", GetUnreadCount(notificationService))
	notifications.Put("/:id/read", MarkNotificationAsRead(notificationService))
	notifications.Put("/read-all", MarkAllNotificationsAsRead(notificationService))
}
