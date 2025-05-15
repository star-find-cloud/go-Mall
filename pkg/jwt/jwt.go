package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/star-find-cloud/star-mall/conf"
	"time"
)

const (
	TokenExpireDuration = 6 * time.Hour
)

var secretKey = conf.GetConfig().App.JWTSecret

// 生成 jwt
func GenerateToken(userID int64, username string) (string, error) {
	// 创建自定义的 claims 对象
	claims := CustomClaims{
		UserID:   userID,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "star-Mall",
		},
	}
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 返回签名后的token
	return token.SignedString([]byte(secretKey))
}

// 解析token
func ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
