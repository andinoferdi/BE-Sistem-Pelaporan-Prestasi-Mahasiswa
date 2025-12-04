package service

import (
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repository "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	utilspostgre "sistem-pelaporan-prestasi-mahasiswa/utils/postgre"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoginService(c *fiber.Ctx, db *sql.DB) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
			},
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Username dan password wajib diisi.",
			},
		})
	}

	user, err := repository.GetUserByUsernameOrEmail(db, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Username atau password tidak valid.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data user dari database. Detail: " + err.Error(),
			},
		})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akun Anda tidak aktif. Silakan hubungi administrator.",
			},
		})
	}

	if !utilspostgre.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Username atau password tidak valid.",
			},
		})
	}

	token, err := utilspostgre.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error generating token. Detail: " + err.Error(),
			},
		})
	}

	refreshToken, err := utilspostgre.GenerateRefreshToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error generating refresh token. Detail: " + err.Error(),
			},
		})
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	if err := repository.SaveRefreshToken(db, user.ID, refreshToken, expiresAt); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error menyimpan refresh token. Detail: " + err.Error(),
			},
		})
	}

	permissions, err := repository.GetUserPermissions(db, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil permissions. Detail: " + err.Error(),
			},
		})
	}

	roleName, err := repository.GetRoleName(db, user.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil role name. Detail: " + err.Error(),
			},
		})
	}

	response := model.LoginResponse{
		Status: "success",
		Data: struct {
			Token        string           `json:"token"`
			RefreshToken string           `json:"refreshToken"`
			User         model.LoginUserResponse `json:"user"`
		}{
			Token:        token,
			RefreshToken: refreshToken,
			User: model.LoginUserResponse{
				ID:          user.ID,
				Username:    user.Username,
				FullName:    user.FullName,
				Role:        roleName,
				Permissions: permissions,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetProfileService(c *fiber.Ctx, db *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Data user tidak ditemukan di database.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data user dari database. Detail: " + err.Error(),
			},
		})
	}

	response := model.GetProfileResponse{
		Status: "success",
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
	var req model.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
			},
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Refresh token wajib diisi.",
			},
		})
	}

	_, err := repository.GetRefreshToken(db, req.RefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Refresh token tidak valid atau sudah expired.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil refresh token dari database. Detail: " + err.Error(),
			},
		})
	}

	claims, err := utilspostgre.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		repository.DeleteRefreshToken(db, req.RefreshToken)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Refresh token tidak valid atau sudah expired.",
			},
		})
	}

	user, err := repository.GetUserByID(db, claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Data user tidak ditemukan di database.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data user dari database. Detail: " + err.Error(),
			},
		})
	}

	if !user.IsActive {
		repository.DeleteRefreshToken(db, req.RefreshToken)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akun Anda tidak aktif. Silakan hubungi administrator.",
			},
		})
	}

	token, err := utilspostgre.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error generating token. Detail: " + err.Error(),
			},
		})
	}

	refreshToken, err := utilspostgre.GenerateRefreshToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error generating refresh token. Detail: " + err.Error(),
			},
		})
	}

	repository.DeleteRefreshToken(db, req.RefreshToken)

	expiresAt := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	if err := repository.SaveRefreshToken(db, user.ID, refreshToken, expiresAt); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error menyimpan refresh token. Detail: " + err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        token,
			"refreshToken": refreshToken,
		},
	})
}

func LogoutService(c *fiber.Ctx, db *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	if err := repository.DeleteUserRefreshTokens(db, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error menghapus refresh token. Detail: " + err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

func HealthCheckService(c *fiber.Ctx) error {
	instanceID := c.Locals("server_instance_id")
	if instanceID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Server instance ID tidak ditemukan.",
			},
		})
	}

	response := fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"instanceId": instanceID,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

