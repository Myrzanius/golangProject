package config

import (
	"log"
	"os"
)

var (
	InventoryServiceURL string
	OrderServiceURL     string
	AuthServiceURL      string
)

func init() {
	InventoryServiceURL = getEnv("INVENTORY_SERVICE_URL", "http://localhost:8081")
	OrderServiceURL = getEnv("ORDER_SERVICE_URL", "http://localhost:8082")
	AuthServiceURL = getEnv("AUTH_SERVICE_URL", "http://localhost:8083")
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Warning: environment variable %s not set, using default value: %s", key, fallback)
		return fallback
	}
	return value
}
