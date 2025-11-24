package middleware

import (
	utilspostgre "sistem-pelaporan-prestasi-mahasiswa/utils/postgre"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token akses diperlukan. Tambahkan header 'Authorization: Bearer YOUR_TOKEN'.",
			})
		}

		tokenString := utilspostgre.ExtractTokenFromHeader(authHeader)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Format token tidak valid. Gunakan format 'Bearer YOUR_TOKEN'.",
			})
		}

		claims, err := utilspostgre.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid atau sudah expired. Silakan login ulang untuk mendapatkan token baru.",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role_id", claims.RoleID)

		return c.Next()
	}
}

func RoleRequired(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Role tidak ditemukan. Silakan login ulang.",
			})
		}

		for _, role := range allowedRoles {
			if roleID == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Akses ditolak. Anda tidak memiliki permission untuk mengakses endpoint ini.",
		})
	}
}

func PermissionRequired(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		hasPermission, err := utilspostgre.CheckUserPermission(userID, permission)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal memeriksa permission: " + err.Error(),
			})
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak. Anda tidak memiliki permission '" + permission + "'.",
			})
		}

		return c.Next()
	}
}

