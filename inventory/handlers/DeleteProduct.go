package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"inventory/models"
	"net/http"
)

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID to ObjectID type
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Check if product exists in products collection
	var product models.Product
	err = productsCollection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to check if product exists", http.StatusInternalServerError)
		}
		return
	}

	// Check if product is referenced in any order_items
	var orderCount int64
	orderFilter := bson.M{"product_id": productID}
	orderCount, err = ordersCollection.CountDocuments(context.Background(), orderFilter)
	if err != nil {
		http.Error(w, "Failed to check orders referencing product", http.StatusInternalServerError)
		return
	}
	if orderCount > 0 {
		http.Error(w, "Cannot delete product that is referenced in orders", http.StatusConflict)
		return
	}

	// Delete the product from the products collection
	_, err = productsCollection.DeleteOne(context.Background(), bson.M{"_id": productID})
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	// Respond with HTTP 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
