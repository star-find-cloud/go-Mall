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

func InitRouter(userHandler *handler.UserHandler, imageHandler *handler.ImageHandler) *gin.Engine {
	// 设置 gin 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(
		middleware.GinLogger(),
		middleware.GinRecoveryWithZap(true),
	)

	publicGroup := r.Group("/api")
	{
		publicGroup.POST("/login", userHandler.Login)
		publicGroup.POST("/register", userHandler.Register)
		publicGroup.GET("/health", healthCheck)
	}

	authGroup := r.Group("/api")
	authGroup.Use(middleware.JwtAuth())
	{
		//authGroup.POST("/update", userHandler.)
	}

	imageGroup := r.Group("/image")
	{
		imageGroup.POST("/upload", imageHandler.UploadImage)
		imageGroup.POST("/uploads", imageHandler.UploadImages)
	}

	return r
}
