package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	appjwt "github.com/star-find-cloud/star-mall/pkg/jwt"
	"net/http"
)

func handleJWTError(c *gin.Context, err error) {
	if errors.Is(err, jwt.ErrTokenExpired) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "token is expired",
		})
	} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "token is not active yet",
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "token is invalid",
		})
	}
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "token is empty",
			})
			return
		}

		const bearerPrefix = "Bearer "
		tokenStr := authHeader
		if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
			tokenStr = authHeader[len(bearerPrefix):]
		}

		claims, err := appjwt.ParseToken(tokenStr)
		if err != nil {
			handleJWTError(c, err)
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
