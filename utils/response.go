package utils

import (
	"github.com/gin-gonic/gin"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
)

type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

// 返回 json 格式的 响应
func RespondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"code":    status,
		"data":    data,
		"message": "success",
	})
}

// 返回 protobuf 格式的 响应
func RespondProtobuf(c *gin.Context, status int, data interface{}) {
	c.ProtoBuf(status, gin.H{
		"code":    status,
		"data":    data,
		"message": "success",
	})
}

func RespondError(c *gin.Context, status int, message string, err interface{}) {
	applog.AppLogger.Errorf("error: %v", err)
	c.AbortWithStatusJSON(status, ResponseError{
		Code:    status,
		Message: message,
		Error:   err,
	})
}
