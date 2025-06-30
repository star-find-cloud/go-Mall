package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
	"runtime/debug"
	"time"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.RawPath
		rawQuery := c.Request.URL.RawQuery

		c.Next()
		// 获取耗时
		latency := time.Since(start)
		// 定义 httpLogger 信息结构
		log.HttpLogger.Infow(
			"Request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", rawQuery,
			"ip", c.ClientIP(),
			"user-agent", c.Request.UserAgent(),
			"latency", latency.String(),
			"errors", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}

func GinRecoveryWithZap(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 网络错误处理
				if utils.IsBrokenPipeErr(err) {
					log.AppLogger.Errorw("Broken pipe",
						"error", err,
						"request", utils.DumpRequest(c),
					)
					c.Abort()
					return
				}
				// 带堆栈的错误处理
				if stack {
					log.AppLogger.Errorw("[Recovery from panic]",
						"error", err,
						"stack", string(debug.Stack()),
					)
				} else {
					log.AppLogger.Errorw("[Recovery from panic]",
						"error", err,
					)
				}
				// 出现错误时, 返回500错误
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
