package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func NotificationRoutes(app *fiber.App, postgresDB *sql.DB) {
	notifications := app.Group("/api/v1/notifications", middlewarepostgre.AuthRequired())

	notifications.Get("", func(c *fiber.Ctx) error { return servicepostgre.GetNotificationsService(c, postgresDB) })
	notifications.Get("/unread-count", func(c *fiber.Ctx) error { return servicepostgre.GetUnreadCountService(c, postgresDB) })
	notifications.Put("/:id/read", func(c *fiber.Ctx) error { return servicepostgre.MarkAsReadService(c, postgresDB) })
	notifications.Put("/read-all", func(c *fiber.Ctx) error { return servicepostgre.MarkAllAsReadService(c, postgresDB) })
}

