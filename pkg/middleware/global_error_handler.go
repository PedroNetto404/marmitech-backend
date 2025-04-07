package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("🔥 panic recovered: %v", r)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal_server_error",
					"msg":   "Oops! Something went wrong.",
				})
			}
		}()
		c.Next()

		if len(c.Errors) > 0 {
			log.Printf("⚠️ gin error: %v", c.Errors.String())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "unexpected_error",
				"msg":   "An internal error occurred.",
			})
		}
	}
}
