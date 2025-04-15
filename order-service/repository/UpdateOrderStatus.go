package repository

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"inventory/handlers"

	"net/http"
	"time"
)

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
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
	var existingOrder Order
	err = handlers.ordersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingOrder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			handlers.RespondWithError(w, http.StatusNotFound, "Order not found")
			return
		}
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Parse request
	var req handlers.UpdateOrderStatusRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate status
	validStatuses := map[string]bool{"pending": true, "processing": true, "shipped": true, "delivered": true, "cancelled": true}
	if !validStatuses[req.Status] {
		handlers.RespondWithError(w, http.StatusBadRequest, "Invalid status. Must be one of: pending, processing, shipped, delivered, cancelled")
		return
	}

	// Start a session for transaction if we're cancelling an order
	session, err := handlers.Client.StartSession()
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer session.EndSession(ctx)

	// Handle cancellation with transaction
	if existingOrder.Status != "cancelled" && req.Status == "cancelled" {
		err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
			// Transaction started
			if err := session.StartTransaction(); err != nil {
				return err
			}

			// Restore stock for each item
			for _, item := range existingOrder.Items {
				_, err = handlers.productsCollection.UpdateOne(
					sessionContext,
					bson.M{"_id": item.ProductID},
					bson.M{"$inc": bson.M{"stock_level": item.Quantity}, "$set": bson.M{"updated_at": time.Now()}},
				)
				if err != nil {
					return err
				}
			}

			// Update order status
			_, err = handlers.ordersCollection.UpdateOne(
				sessionContext,
				bson.M{"_id": id},
				bson.M{"$set": bson.M{"status": req.Status, "updated_at": time.Now()}},
			)
			if err != nil {
				return err
			}

			// Commit the transaction
			return session.CommitTransaction(sessionContext)
		})

		if err != nil {
			// Transaction failed, handle the error
			session.AbortTransaction(ctx)
			handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// Simple status update (no stock restoration needed)
		_, err = handlers.ordersCollection.UpdateOne(
			ctx,
			bson.M{"_id": id},
			bson.M{"$set": bson.M{"status": req.Status, "updated_at": time.Now()}},
		)
		if err != nil {
			handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Get updated order
	var updatedOrder Order
	err = handlers.ordersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedOrder)
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, updatedOrder)
}
