package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/star-find-cloud/star-mall/handler"
	"github.com/star-find-cloud/star-mall/middleware"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
)

type healthCheckResponse struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

// 健康检查
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, healthCheckResponse{
		Status:   "UP",
		Hostname: utils.GetHostName(),
	})
}

func InitRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.GinLogger(),
		middleware.GinRecoveryWithZap(true),
	)

	authGroup := r.Group("/api")
	authGroup.Use(middleware.JwtAuth())
	{
		authGroup.GET("/login", userHandler.Login)
	}

	return r
}
