package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"inventory/models"
	"inventory/repository"
	"net/http"
	"time"
)

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert string ID to ObjectID
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		repository.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Check if product exists
	var existingProduct models.Product
	err = productsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			repository.RespondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Parse update fields
	var updateFields map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updateFields)
	if err != nil {
		repository.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Build update document
	update := bson.M{"updated_at": time.Now()}

	// Handle each field
	if name, ok := updateFields["name"].(string); ok && name != "" {
		update["name"] = name
	}
	if description, ok := updateFields["description"].(string); ok {
		update["description"] = description
	}
	if price, ok := updateFields["price"].(float64); ok && price > 0 {
		update["price"] = price
	}
	if stockLevel, ok := updateFields["stock_level"].(float64); ok && stockLevel >= 0 {
		update["stock_level"] = int(stockLevel)
	}
	if categoryIDStr, ok := updateFields["category_id"].(string); ok {
		// Convert string category ID to ObjectID
		categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
		if err != nil {
			repository.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
			return
		}

		// Check if category exists
		var categoryDoc models.Category
		err = categoriesCollection.FindOne(ctx, bson.M{"_id": categoryID}).Decode(&categoryDoc)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				repository.RespondWithError(w, http.StatusBadRequest, "Category not found")
				return
			}
			repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		update["category_id"] = categoryID
	}

	// If no valid fields to update
	if len(update) == 1 { // Only updated_at is set
		repository.RespondWithError(w, http.StatusBadRequest, "No valid fields to update")
		return
	}

	// Update product
	result, err := productsCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)
	if err != nil {
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if result.ModifiedCount == 0 {
		repository.RespondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	// Get updated product
	var updatedProduct models.Product
	err = productsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedProduct)
	if err != nil {
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	repository.RespondWithJSON(w, http.StatusOK, updatedProduct)
}
