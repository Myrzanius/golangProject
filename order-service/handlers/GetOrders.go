package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"inventory/handlers"
	"net/http"
	"strconv"
)

func GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get query parameters for filtering
	userIDStr := r.URL.Query().Get("user_id")
	status := r.URL.Query().Get("status")

	// Build the filter
	filter := bson.M{}

	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err == nil {
			filter["user_id"] = userID
		}
	}

	if status != "" {
		filter["status"] = status
	}

	// Query orders
	cursor, err := handlers.ordersCollection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer cursor.Close(ctx)

	// Decode results
	var orders []Order
	if err := cursor.All(ctx, &orders); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, orders)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert string ID to ObjectID
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	// Find order by ID
	var order Order
	err = handlers.ordersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			handlers.RespondWithError(w, http.StatusNotFound, "Order not found")
			return
		}
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, order)
}
