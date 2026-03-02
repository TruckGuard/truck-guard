package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)
		method := c.Request.Method
		clientIP := c.ClientIP()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				slog.Error("request error",
					"method", method,
					"path", path,
					"query", query,
					"ip", clientIP,
					"status", status,
					"latency", latency,
					"error", e,
				)
			}
		} else {
			slog.Info("request",
				"method", method,
				"path", path,
				"query", query,
				"ip", clientIP,
				"status", status,
				"latency", latency,
			)
		}
	}
}
