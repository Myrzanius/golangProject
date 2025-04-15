package middleware

import (
	"log"
	"time"
)

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		log.Printf("Started %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("Completed %s %s with status %d in %v", c.Request.Method, c.Request.URL.Path, statusCode, duration)
	}
}
