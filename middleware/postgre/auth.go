package middleware

// #1 proses: import library yang diperlukan untuk database, utils, dan fiber
import (
	"database/sql"
	utilspostgre "sistem-pelaporan-prestasi-mahasiswa/utils/postgre"

	"github.com/gofiber/fiber/v2"
)

// #2 proses: middleware untuk validasi JWT token dan set user info ke context
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// #2a proses: ambil Authorization header dari request
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Token akses diperlukan. Tambahkan header 'Authorization: Bearer YOUR_TOKEN'.",
			})
		}

		// #2b proses: extract token dari header Authorization
		tokenString := utilspostgre.ExtractTokenFromHeader(authHeader)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Format token tidak valid. Gunakan format 'Bearer YOUR_TOKEN'.",
			})
		}

		// #2c proses: validasi token dan ambil claims
		claims, err := utilspostgre.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Token tidak valid atau sudah expired. Silakan login ulang untuk mendapatkan token baru.",
			})
		}

		// #2d proses: set user info ke context untuk digunakan di handler
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role_id", claims.RoleID)

		return c.Next()
	}
}

// #3 proses: middleware untuk validasi role user, hanya allow role yang diizinkan
func RoleRequired(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// #3a proses: ambil role ID dari context yang diset oleh AuthRequired
		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "Role tidak ditemukan. Silakan login ulang.",
			})
		}

		// #3b proses: cek apakah role ID user ada di list allowed roles
		for _, role := range allowedRoles {
			if roleID == role {
				return c.Next()
			}
		}

		// #3c proses: jika role tidak diizinkan, return forbidden response
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Akses ditolak",
			"message": "Akses ditolak. Anda tidak memiliki permission untuk mengakses endpoint ini.",
		})
	}
}

// #4 proses: middleware untuk validasi permission user, cek di database
func PermissionRequired(db *sql.DB, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// #4a proses: ambil user ID dari context yang diset oleh AuthRequired
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Tidak diizinkan",
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			})
		}

		// #4b proses: cek apakah user memiliki permission tertentu di database
		hasPermission, err := utilspostgre.CheckUserPermission(db, userID, permission)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": "Gagal memeriksa permission: " + err.Error(),
			})
		}

		// #4c proses: jika user tidak punya permission, return forbidden response
		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Akses ditolak",
				"message": "Akses ditolak. Anda tidak memiliki permission '" + permission + "'.",
			})
		}

		return c.Next()
	}
}
