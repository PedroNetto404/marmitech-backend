package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func HttpLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		duration := time.Since(start)

		status := c.Writer.Status()
		log.Printf(
			"%s | %3d | %13v | %-15s |%-7s  %s",
			time.Now().Format(time.RFC3339),
			status,
			duration,
			c.ClientIP(),
			c.Request.Method,
			path,
		)
	}
}
