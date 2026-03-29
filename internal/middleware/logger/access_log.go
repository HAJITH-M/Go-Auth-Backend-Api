package middleware

import (
	"go-auth-backend-api/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Init()
		start := time.Now()
		c.Next()

		logger.Log.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("duration", time.Since(start)),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("remote_addr", c.ClientIP()),
			zap.Int("status", c.Writer.Status()),
		)
	}
}
