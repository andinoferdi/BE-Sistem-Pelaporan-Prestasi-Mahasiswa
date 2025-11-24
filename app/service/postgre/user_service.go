package service

import (
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repository "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	utilspostgre "sistem-pelaporan-prestasi-mahasiswa/utils/postgre"

	"github.com/gofiber/fiber/v2"
)

func LoginService(c *fiber.Ctx, db *sql.DB) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email dan password wajib diisi.",
		})
	}

	user, err := repository.GetUserByEmail(db, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Email atau password tidak valid.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data user dari database. Detail: " + err.Error(),
		})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Akun Anda tidak aktif. Silakan hubungi administrator.",
		})
	}

	if !utilspostgre.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Email atau password tidak valid.",
		})
	}

	token, err := utilspostgre.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error generating token. Detail: " + err.Error(),
		})
	}

	user.PasswordHash = ""

	response := model.LoginResponse{
		Success: true,
		Message: "Login berhasil.",
		Data: struct {
			User  model.User `json:"user"`
			Token string     `json:"token"`
		}{
			User:  *user,
			Token: token,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetProfileService(c *fiber.Ctx, db *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User ID tidak ditemukan. Silakan login ulang.",
		})
	}

	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data user tidak ditemukan di database.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data user dari database. Detail: " + err.Error(),
		})
	}

	response := model.GetProfileResponse{
		Success: true,
		Message: "Data profile berhasil diambil.",
		Data: struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			FullName string `json:"full_name"`
			RoleID   string `json:"role_id"`
		}{
			UserID:   user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			RoleID:   user.RoleID,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func RefreshTokenService(c *fiber.Ctx, db *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User ID tidak ditemukan. Silakan login ulang.",
		})
	}

	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data user tidak ditemukan di database.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error mengambil data user dari database. Detail: " + err.Error(),
		})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Akun Anda tidak aktif. Silakan hubungi administrator.",
		})
	}

	token, err := utilspostgre.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error generating token. Detail: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Token berhasil di-refresh.",
		"data": fiber.Map{
			"token": token,
		},
	})
}

func LogoutService(c *fiber.Ctx, db *sql.DB) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logout berhasil.",
	})
}

