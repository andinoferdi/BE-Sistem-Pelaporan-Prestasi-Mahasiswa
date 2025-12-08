package route

import (
	"database/sql"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, db *sql.DB) {
	users := app.Group("/api/v1/users", middlewarepostgre.AuthRequired())

	users.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.GetAllUsersService(c, db)
	})

	users.Get("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.GetUserByIDService(c, db)
	})

	users.Post("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.CreateUserService(c, db)
	})

	users.Put("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.UpdateUserService(c, db)
	})

	users.Delete("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.DeleteUserService(c, db)
	})

	users.Put("/:id/role", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		return servicepostgre.UpdateUserRoleService(c, db)
	})
}

