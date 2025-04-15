package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"inventory/models"
	"inventory/repository"
	"net/http"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Query categories
	cursor, err := categoriesCollection.Find(ctx, bson.M{})
	if err != nil {
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer cursor.Close(ctx)

	// Decode results
	var categories []models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	repository.RespondWithJSON(w, http.StatusOK, categories)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		repository.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Basic validation
	if category.Name == "" {
		repository.RespondWithError(w, http.StatusBadRequest, "Category name is required")
		return
	}

	// Check if category name already exists
	count, err := categoriesCollection.CountDocuments(ctx, bson.M{"name": category.Name})
	if err != nil {
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if count > 0 {
		repository.RespondWithError(w, http.StatusConflict, "Category name already exists")
		return
	}

	// Create new category with new ID
	category.ID = primitive.NewObjectID()

	// Insert category
	_, err = categoriesCollection.InsertOne(ctx, category)
	if err != nil {
		repository.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	repository.RespondWithJSON(w, http.StatusCreated, category)
}
