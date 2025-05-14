package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "token is empty",
			})
			return
		}

		claims, err := ParseToken(tokenStr)
		if err != nil {
			handleJWTError(c, err)
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
