package utils

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

func HashPassword(password string) (string, error) {
	// 获取被哈希加密后的 password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	return string(bytes), err
}

// CheckPasswordHash 检查密码是否正确
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// VerifyEmail 验证邮箱
func VerifyEmail(email string) bool {
	reg := regexp.MustCompile(`^[A-Za-z0-9\p{Han}]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	return reg.MatchString(email)
}

// Bcrypt 生成 Bcrypt 加密
func Bcrypt(str string) (string, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	err := bcrypt.CompareHashAndPassword(hash, []byte(str))

	return string(hash), err
}
