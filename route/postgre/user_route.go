package route

// #1 proses: import library yang diperlukan untuk context, database, model, service, middleware, strings, time, dan fiber
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

// GetAllUsers godoc
// @Summary Get all users
// @Description Mengambil daftar semua user dari database. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} model.GetAllUsersResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users [get]
func GetAllUsers(userService servicepostgre.IUserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Mengambil data user spesifik berdasarkan ID. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} model.GetUserByIDResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/{id} [get]
func GetUserByID(userService servicepostgre.IUserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// CreateUser godoc
// @Summary Create user
// @Description Membuat user baru di database. Dapat sekaligus membuat student atau lecturer profile. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body model.CreateUserRequest true "User data"
// @Success 200 {object} model.CreateUserResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 409 {object} map[string]string "Conflict"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /users [post]
func CreateUser(userService servicepostgre.IUserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// UpdateUser godoc
// @Summary Update user
// @Description Memperbarui data user berdasarkan ID. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param body body model.UpdateUserRequest true "User data"
// @Success 200 {object} model.UpdateUserResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 409 {object} map[string]string "Conflict"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /users/{id} [put]
func UpdateUser(userService servicepostgre.IUserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// DeleteUser godoc
// @Summary Delete user
// @Description Menghapus user berdasarkan ID. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} model.DeleteUserResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 409 {object} map[string]string "Conflict"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/{id} [delete]
func DeleteUser(userService servicepostgre.IUserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// UpdateUserRole godoc
// @Summary Update user role
// @Description Memperbarui role user berdasarkan ID. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param body body object true "Role ID" example({"role_id": "uuid-here"})
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 409 {object} map[string]string "Conflict"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /users/{id}/role [put]
func UpdateUserRole(userService servicepostgre.IUserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// CreateStudentProfile godoc
// @Summary Create student profile
// @Description Membuat student profile untuk user tertentu. User harus memiliki role Mahasiswa. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param body body model.CreateStudentRequest true "Student data"
// @Success 200 {object} model.CreateStudentResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 409 {object} map[string]string "Conflict"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /users/{id}/student-profile [post]
func CreateStudentProfile(studentService servicepostgre.IStudentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// CreateLecturerProfile godoc
// @Summary Create lecturer profile
// @Description Membuat lecturer profile untuk user tertentu. User harus memiliki role Dosen Wali. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param body body model.CreateLecturerRequest true "Lecturer data"
// @Success 200 {object} model.CreateLecturerResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 409 {object} map[string]string "Conflict"
// @Failure 422 {object} map[string]string "Unprocessable Entity"
// @Router /users/{id}/lecturer-profile [post]
func CreateLecturerProfile(lecturerService servicepostgre.ILecturerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// GetAllRoles godoc
// @Summary Get all roles
// @Description Mengambil daftar semua role dari database. Hanya dapat diakses oleh admin dengan permission user:manage
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} model.GetAllRolesResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /roles [get]
func GetAllRoles(userService servicepostgre.IUserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// #2 proses: setup semua route untuk user dengan middleware AuthRequired dan PermissionRequired
func UserRoutes(app *fiber.App, userService servicepostgre.IUserService, studentService servicepostgre.IStudentService, lecturerService servicepostgre.ILecturerService, db *sql.DB) {
	users := app.Group("/api/v1/users", middlewarepostgre.AuthRequired())

	users.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), GetAllUsers(userService))

	users.Get("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), GetUserByID(userService))

	users.Post("", middlewarepostgre.PermissionRequired(db, "user:manage"), CreateUser(userService))

	users.Put("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), UpdateUser(userService))

	users.Delete("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), DeleteUser(userService))

	users.Put("/:id/role", middlewarepostgre.PermissionRequired(db, "user:manage"), UpdateUserRole(userService))

	users.Post("/:id/student-profile", middlewarepostgre.PermissionRequired(db, "user:manage"), CreateStudentProfile(studentService))

	users.Post("/:id/lecturer-profile", middlewarepostgre.PermissionRequired(db, "user:manage"), CreateLecturerProfile(lecturerService))

	roles := app.Group("/api/v1/roles", middlewarepostgre.AuthRequired())

	roles.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), GetAllRoles(userService))
}
