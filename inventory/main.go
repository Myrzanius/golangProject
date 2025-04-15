// main.go
package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"inventory/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	handlers.InitMongo(ctx)
	defer handlers.Client.Disconnect(ctx)

	r := mux.NewRouter()

	// Product endpoints
	r.HandleFunc("/products", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handlers.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PATCH")
	r.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE")

	// Category endpoints
	r.HandleFunc("/categories", handlers.GetCategories).Methods("GET")
	r.HandleFunc("/categories", handlers.CreateCategory).Methods("POST")

	// Order endpoints
	r.HandleFunc("/orders", handlers.GetOrders).Methods("GET")
	r.HandleFunc("/orders", handlers.CreateOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", handlers.GetOrder).Methods("GET")
	r.HandleFunc("/orders/{id}", handlers.UpdateOrderStatus).Methods("PATCH")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
