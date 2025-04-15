package middleware

import (
	"fmt"
	"log"
	"time"
)

// LoggingMiddleware - renamed to match usage in main.go (capital L)
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		log.Printf("Started %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("Completed %s %s with status %d in %v", c.Request.Method, c.Request.URL.Path, statusCode, duration)
	}
}

// TelemetryMiddleware - renamed to match usage in main.go (capital T)
func TelemetryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// Start timer for request duration
		startTime := time.Now()

		// Start span for OpenTelemetry
		tr := otel.Tracer("api-gateway")
		ctx, span := tr.Start(c.Request.Context(), fmt.Sprintf("%s %s", method, path))
		defer span.End()

		// Store context with span in request
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Record metrics after completion
		duration := time.Since(startTime).Seconds()
		status := fmt.Sprintf("%d", c.Writer.Status())
		service := "unknown"
		endpoint := path

		if len(c.Params) > 0 {
			// Extract service and endpoint info from path
			if c.FullPath() != "" {
				endpoint = c.FullPath()
			}
		}

		// Determine service based on path
		if len(path) > 1 {
			if path[:9] == "/products" {
				service = "inventory"
			} else if path[:7] == "/orders" {
				service = "orders"
			}
		}

		// Add attributes to the span
		span.SetAttributes(
			attribute.String("http.method", method),
			attribute.String("http.path", path),
			attribute.Int("http.status_code", c.Writer.Status()),
		)
	}
}
