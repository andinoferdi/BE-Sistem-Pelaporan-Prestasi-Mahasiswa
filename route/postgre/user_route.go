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

// #2 proses: setup semua route untuk user dengan middleware AuthRequired dan PermissionRequired
func UserRoutes(app *fiber.App, userService servicepostgre.IUserService, studentService servicepostgre.IStudentService, lecturerService servicepostgre.ILecturerService, db *sql.DB) {
	// #2a proses: buat group route untuk users dengan middleware AuthRequired
	users := app.Group("/api/v1/users", middlewarepostgre.AuthRequired())

	// #3 proses: endpoint GET /api/v1/users untuk ambil semua user dengan permission user:manage
	users.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #3a proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #3b proses: panggil service get all users
		usersList, err := userService.GetAllUsers(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #3c proses: build response dengan data users
		response := model.GetAllUsersResponse{
			Status: "success",
			Data:   usersList,
		}

		return c.JSON(response)
	})

	// #4 proses: endpoint GET /api/v1/users/:id untuk ambil user berdasarkan ID dengan permission user:manage
	users.Get("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #4a proses: ambil user ID dari URL parameter dan validasi
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		// #4b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #4c proses: panggil service get user by ID
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

		// #4d proses: build response dengan data user
		response := model.GetUserByIDResponse{
			Status: "success",
			Data:   *user,
		}

		return c.JSON(response)
	})

	// #5 proses: endpoint POST /api/v1/users untuk buat user baru dengan permission user:manage
	users.Post("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #5a proses: parse request body ke CreateUserRequest
		req := new(model.CreateUserRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #5b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #5c proses: panggil service create user
		user, err := userService.CreateUser(ctx, *req)
		if err != nil {
			// #5d proses: handle error dengan status code yang sesuai
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

		// #5e proses: build response dengan data user yang baru dibuat
		response := model.CreateUserResponse{
			Status: "success",
			Data:   *user,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #6 proses: endpoint PUT /api/v1/users/:id untuk update user dengan permission user:manage
	users.Put("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #6a proses: ambil user ID dari URL parameter dan validasi
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		// #6b proses: parse request body ke UpdateUserRequest
		req := new(model.UpdateUserRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #6c proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #6d proses: panggil service update user
		user, err := userService.UpdateUser(ctx, id, *req)
		if err != nil {
			// #6e proses: handle error dengan status code yang sesuai
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

		// #6f proses: build response dengan data user yang sudah diupdate
		response := model.UpdateUserResponse{
			Status: "success",
			Data:   *user,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #7 proses: endpoint DELETE /api/v1/users/:id untuk hapus user dengan permission user:manage
	users.Delete("/:id", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #7a proses: ambil user ID dari URL parameter dan validasi
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		// #7b proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #7c proses: panggil service delete user
		err := userService.DeleteUser(ctx, id)
		if err != nil {
			// #7d proses: handle error dengan status code yang sesuai
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

		// #7e proses: build response sukses
		response := model.DeleteUserResponse{
			Status: "success",
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #8 proses: endpoint PUT /api/v1/users/:id/role untuk update role user dengan permission user:manage
	users.Put("/:id/role", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #8a proses: ambil user ID dari URL parameter dan validasi
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		// #8b proses: parse request body untuk ambil role ID
		req := struct {
			RoleID string `json:"role_id"`
		}{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #8c proses: validasi role ID tidak kosong
		if req.RoleID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "role_id wajib diisi.",
			})
		}

		// #8d proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #8e proses: panggil service update user role
		err := userService.UpdateUserRole(ctx, id, req.RoleID)
		if err != nil {
			// #8f proses: handle error dengan status code yang sesuai
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

		// #8g proses: return response sukses
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Role user berhasil diupdate",
		})
	})

	// #9 proses: endpoint POST /api/v1/users/:id/student-profile untuk buat student profile untuk user dengan permission user:manage
	users.Post("/:id/student-profile", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #9a proses: ambil user ID dari URL parameter dan validasi
		userID := c.Params("id")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		// #9b proses: parse request body ke CreateStudentRequest
		req := new(model.CreateStudentRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #9c proses: set user ID ke request
		req.UserID = userID

		// #9d proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #9e proses: panggil service create student
		student, err := studentService.CreateStudent(ctx, *req)
		if err != nil {
			// #9f proses: handle error dengan status code yang sesuai
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

		// #9g proses: build response dengan data student yang baru dibuat
		response := model.CreateStudentResponse{
			Status: "success",
			Data:   *student,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #10 proses: endpoint POST /api/v1/users/:id/lecturer-profile untuk buat lecturer profile untuk user dengan permission user:manage
	users.Post("/:id/lecturer-profile", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #10a proses: ambil user ID dari URL parameter dan validasi
		userID := c.Params("id")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "ID user wajib diisi.",
			})
		}

		// #10b proses: parse request body ke CreateLecturerRequest
		req := new(model.CreateLecturerRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		// #10c proses: set user ID ke request
		req.UserID = userID

		// #10d proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #10e proses: panggil service create lecturer
		lecturer, err := lecturerService.CreateLecturer(ctx, *req)
		if err != nil {
			// #10f proses: handle error dengan status code yang sesuai
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

		// #10g proses: build response dengan data lecturer yang baru dibuat
		response := model.CreateLecturerResponse{
			Status: "success",
			Data:   *lecturer,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	// #11 proses: buat group route untuk roles dengan middleware AuthRequired
	roles := app.Group("/api/v1/roles", middlewarepostgre.AuthRequired())

	// #12 proses: endpoint GET /api/v1/roles untuk ambil semua role dengan permission user:manage
	roles.Get("", middlewarepostgre.PermissionRequired(db, "user:manage"), func(c *fiber.Ctx) error {
		// #12a proses: buat context dengan timeout 5 detik
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// #12b proses: panggil service get all roles
		rolesList, err := userService.GetAllRoles(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}

		// #12c proses: build response dengan data roles
		response := model.GetAllRolesResponse{
			Status: "success",
			Data:   rolesList,
		}

		return c.JSON(response)
	})
}
