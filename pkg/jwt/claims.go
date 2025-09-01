package jwt

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID   int64
	UserName string
	Roles    int64
	jwt.RegisteredClaims
}
