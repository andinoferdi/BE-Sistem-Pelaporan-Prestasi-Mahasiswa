// #1 proses: package untuk utility functions terkait password hashing
package postgre

// #2 proses: import library bcrypt untuk secure password hashing
import "golang.org/x/crypto/bcrypt"

// #3 proses: hash password menggunakan bcrypt dengan default cost untuk keamanan
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// #4 proses: cek apakah password yang diberikan cocok dengan hash yang tersimpan
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
