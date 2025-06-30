package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// 获取被哈希加密后的 password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	return string(bytes), err
}

// 检查密码是否正确
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
