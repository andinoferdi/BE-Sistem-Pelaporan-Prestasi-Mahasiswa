package route

import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
	middlewarepostgre "sistem-pelaporan-prestasi-mahasiswa/middleware/postgre"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userService servicepostgre.IUserService, studentService servicepostgre.IStudentService, lecturerService servicepostgre.ILecturerService, db *sql.DB) {
	users := app.Group("/api/v1/users", middlewarepostgre.AuthRequired())

	users.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		usersList, err := userService.GetAllUsers(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		response := model.GetAllUsersResponse{
			Status: "success",
			Data:   usersList,
		}

		return c.JSON(response)
	})

	users.Get("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := userService.GetUserByID(ctx, id)
		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Gagal mengambil pengguna",
					"message": err.Error(),
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		response := model.GetUserByIDResponse{
			Status: "success",
			Data:   *user,
		}

		return c.JSON(response)
	})

	users.Post("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		req := new(model.CreateUserRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := userService.CreateUser(ctx, *req)
		if err != nil {
			if strings.Contains(err.Error(), "sudah digunakan") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error":   "Konflik data",
					"message": err.Error(),
				})
			}
			if strings.Contains(err.Error(), "tidak ditemukan") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Data tidak ditemukan",
					"message": err.Error(),
				})
			}
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal membuat user",
				"message": err.Error(),
			})
		}

		response := model.CreateUserResponse{
			Status: "success",
			Data:   *user,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	users.Put("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		req := new(model.UpdateUserRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := userService.UpdateUser(ctx, id, *req)
		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Data tidak ditemukan",
					"message": err.Error(),
				})
			}
			if strings.Contains(err.Error(), "sudah digunakan") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error":   "Konflik data",
					"message": err.Error(),
				})
			}
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal mengupdate user",
				"message": err.Error(),
			})
		}

		response := model.UpdateUserResponse{
			Status: "success",
			Data:   *user,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	users.Delete("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userService.DeleteUser(ctx, id)
		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Data tidak ditemukan",
					"message": err.Error(),
				})
			}
			if strings.Contains(err.Error(), "masih memiliki") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error":   "Konflik data",
					"message": err.Error(),
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal menghapus user",
				"message": err.Error(),
			})
		}

		response := model.DeleteUserResponse{
			Status: "success",
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	users.Put("/:id/role", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		req := struct {
			RoleID string `json:"role_id"`
		}{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		if req.RoleID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "role_id wajib diisi.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userService.UpdateUserRole(ctx, id, req.RoleID)
		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Data tidak ditemukan",
					"message": err.Error(),
				})
			}
			if strings.Contains(err.Error(), "masih memiliki") || strings.Contains(err.Error(), "sudah memiliki") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error":   "Konflik data",
					"message": err.Error(),
				})
			}
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal mengupdate role user",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Role user berhasil diupdate",
		})
	})

	users.Post("/:id/student-profile", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		userID := c.Params("id")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		req := new(model.CreateStudentRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		req.UserID = userID

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		student, err := studentService.CreateStudent(ctx, *req)
		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Data tidak ditemukan",
					"message": err.Error(),
				})
			}
			if strings.Contains(err.Error(), "sudah memiliki") || strings.Contains(err.Error(), "sudah digunakan") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error":   "Konflik data",
					"message": err.Error(),
				})
			}
			if strings.Contains(err.Error(), "harus memiliki role") {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
					"error":   "Validasi gagal",
					"message": err.Error(),
				})
			}
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal membuat student profile",
				"message": err.Error(),
			})
		}

		response := model.CreateStudentResponse{
			Status: "success",
			Data:   *student,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	users.Post("/:id/lecturer-profile", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		userID := c.Params("id")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		req := new(model.CreateLecturerRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		req.UserID = userID

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		lecturer, err := lecturerService.CreateLecturer(ctx, *req)
		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Data tidak ditemukan",
					"message": err.Error(),
				})
			}
			if strings.Contains(err.Error(), "sudah memiliki") || strings.Contains(err.Error(), "sudah digunakan") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error":   "Konflik data",
					"message": err.Error(),
				})
			}
			if strings.Contains(err.Error(), "harus memiliki role") {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
					"error":   "Validasi gagal",
					"message": err.Error(),
				})
			}
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal membuat lecturer profile",
				"message": err.Error(),
			})
		}

		response := model.CreateLecturerResponse{
			Status: "success",
			Data:   *lecturer,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	roles := app.Group("/api/v1/roles", middlewarepostgre.AuthRequired())

	roles.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rolesList, err := userService.GetAllRoles(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		response := model.GetAllRolesResponse{
			Status: "success",
			Data:   rolesList,
		}

		return c.JSON(response)
	})
}
