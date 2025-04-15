package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"inventory/models"
	"inventory/repository"
	"net/http"
	"strconv"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get query parameters for filtering and pagination
	category := r.URL.Query().Get("category")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Default values
	limit := 10
	offset := 0

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Build the filter
	filter := bson.M{}

	if category != "" {
		// First find the category ID
		var categoryDoc models.Category
		err := categoriesCollection.FindOne(ctx, bson.M{"name": primitive.Regex{Pattern: category, Options: "i"}}).Decode(&categoryDoc)
		if err == nil {
			filter["category_id"] = categoryDoc.ID
		}
	}

	// Set up options for pagination
	options := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Query products
	cursor, err := productsCollection.Find(ctx, filter, options)
	if err != nil {
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer cursor.Close(ctx)

	// Decode results
	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	repository.RespondWithJSON(w, http.StatusOK, products)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert string ID to ObjectID
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		repository.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Find product by ID
	var product models.Product
	err = productsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			repository.RespondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	repository.RespondWithJSON(w, http.StatusOK, product)
}
