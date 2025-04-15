package api_gateway

import (
	"inventory/api-gateway/handlers"
	"inventory/api-gateway/middleware"
	"log"
	"os"
	"runtime/trace"
)

func main() {
	r := gin.Default()

	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.TelemetryMiddleware())

	r.GET("/inventory/*action", handlers.ProxyRequestToInventory)
	r.GET("/orders/*action", handlers.ProxyRequestToOrder)

	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting API Gateway: %v", err)
	}
}
