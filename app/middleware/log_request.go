package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LogRoutesMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		userUUID, _ := GetUserUUID(c)
		if userUUID != uuid.Nil {
			log.Info("",
				zap.Any("method", param.Method),
				zap.Any("path", param.Path),
				zap.Any("status", param.StatusCode),
				zap.Any("latency", param.Latency.Round(time.Millisecond).String()),
				zap.Any("user_uuid", userUUID))
		} else {
			log.Info("",
				zap.Any("method", param.Method),
				zap.Any("path", param.Path),
				zap.Any("status", param.StatusCode),
				zap.Any("latency", param.Latency.Round(time.Millisecond).String()))
		}
	}
}
