package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger logs HTTP request and response information.
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		duration := time.Since(start)

		status := c.Writer.Status()
		size := c.Writer.Size()

		if size < 0 {
			size = 0
		}

		logger.Info("HTTP Request",
			zap.String("uri", c.Request.RequestURI),
			zap.String("method", c.Request.Method),
			zap.Duration("duration", duration),
		)

		logger.Info("HTTP Response",
			zap.Int("status", status),
			zap.Int("size", size),
		)

	}
}
