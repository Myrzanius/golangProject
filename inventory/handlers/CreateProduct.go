package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"inventory/models"
	"inventory/repository"
	"net/http"
	"time"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req models.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		repository.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Basic validation
	if req.Name == "" {
		repository.RespondWithError(w, http.StatusBadRequest, "Product name is required")
		return
	}
	if req.Price <= 0 {
		repository.RespondWithError(w, http.StatusBadRequest, "Price must be greater than zero")
		return
	}
	if req.StockLevel < 0 {
		repository.RespondWithError(w, http.StatusBadRequest, "Stock level cannot be negative")
		return
	}

	// Convert string category ID to ObjectID
	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
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

	// Create new product
	now := time.Now()
	product := models.Product{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		StockLevel:  req.StockLevel,
		CategoryID:  categoryID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Insert product
	_, err = productsCollection.InsertOne(ctx, product)
	if err != nil {
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	repository.RespondWithJSON(w, http.StatusCreated, product)
}
