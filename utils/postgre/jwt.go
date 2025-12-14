// #1 proses: package untuk utility functions terkait JWT authentication dan authorization
package postgre

// #2 proses: import library yang diperlukan untuk database, os, model, time, dan JWT
import (
	"database/sql"
	"os"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// #3 proses: struct untuk menyimpan claims JWT yang berisi user ID, email, role ID, dan registered claims
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	RoleID string `json:"role_id"`
	jwt.RegisteredClaims
}

// #4 proses: variable untuk menyimpan JWT secret key yang diambil dari environment atau default
var jwtSecret = []byte(getJWTSecret())

// #5 proses: ambil JWT secret dari environment variable, jika tidak ada gunakan default key
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "sistem-pelaporan-prestasi-mahasiswa-jwt-secret-key-minimum-32-characters-long-for-production-security"
	}
	return secret
}

// #6 proses: generate access token JWT untuk user dengan expiry 24 jam
func GenerateToken(user model.User) (string, error) {
	// #6a proses: buat claims JWT dengan user ID, email, role ID, dan registered claims
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RoleID: user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sistem-pelaporan-prestasi-mahasiswa-api",
			Subject:   "user-authentication",
		},
	}

	// #6b proses: buat token dengan method HS256 dan sign dengan secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// #7 proses: validasi token JWT dan return claims jika token valid
func ValidateToken(tokenString string) (*JWTClaims, error) {
	// #7a proses: parse token dengan claims dan validasi signing method
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// #7b proses: cek apakah token valid dan return claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// #8 proses: extract token dari Authorization header yang berformat "Bearer <token>"
func ExtractTokenFromHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

// #9 proses: generate refresh token JWT untuk user dengan expiry 7 hari
func GenerateRefreshToken(user model.User) (string, error) {
	// #9a proses: buat claims JWT dengan user ID, email, role ID, dan registered claims untuk refresh token
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RoleID: user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sistem-pelaporan-prestasi-mahasiswa-api",
			Subject:   "refresh-token",
		},
	}

	// #9b proses: buat token dengan method HS256 dan sign dengan secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// #10 proses: validasi refresh token menggunakan fungsi ValidateToken yang sama
func ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return ValidateToken(tokenString)
}

// #11 proses: cek apakah user memiliki permission tertentu dengan query ke database
func CheckUserPermission(db *sql.DB, userID string, permission string) (bool, error) {
	// #11a proses: query untuk cek permission user melalui role_permissions dan permissions table
	query := `
		SELECT COUNT(*) > 0
		FROM role_permissions rp
		INNER JOIN permissions p ON rp.permission_id = p.id
		INNER JOIN users u ON u.role_id = rp.role_id
		WHERE u.id = $1 AND p.name = $2
	`

	// #11b proses: execute query dan scan hasil ke variable hasPermission
	var hasPermission bool
	err := db.QueryRow(query, userID, permission).Scan(&hasPermission)
	if err != nil {
		return false, err
	}

	return hasPermission, nil
}
