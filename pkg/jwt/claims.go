package jwt

import jwt "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID   int64
	UserName string
	jwt.RegisteredClaims
}
