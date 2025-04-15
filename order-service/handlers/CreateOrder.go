package handlers

import (
	""
	"context"
	"encoding/json"
	"fmt"
	"github.com/"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	handlers2 "inventory/handlers"
	"inventory/models"
	"net/http"
	"time"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req handlers2.CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handlers2.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request
	if req.UserID <= 0 {
		handlers2.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	if len(req.Items) == 0 {
		handlers2.RespondWithError(w, http.StatusBadRequest, "Order must contain at least one item")
		return
	}

	// Start a session for transaction
	session, err := handlers2.Client.StartSession()
	if err != nil {
		handlers2.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer session.EndSession(ctx)

	// Process the order in a transaction
	var order Order
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// Transaction started
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Calculate total and check inventory
		var total float64
		var orderItems []OrderItem

		for _, item := range req.Items {
			if item.Quantity <= 0 {
				return fmt.Errorf("item quantity must be greater than zero")
			}

			// Convert string product ID to ObjectID
			productID, err := primitive.ObjectIDFromHex(item.ProductID)
			if err != nil {
				return fmt.Errorf("invalid product ID: %s", item.ProductID)
			}

			// Find product
			var product models.Product
			err = handlers2.productsCollection.FindOne(sessionContext, bson.M{"_id": productID}).Decode(&product)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return fmt.Errorf("product with ID %s not found", item.ProductID)
				}
				return err
			}

			// Check stock
			if product.StockLevel < item.Quantity {
				return fmt.Errorf("not enough stock for product with ID %s", item.ProductID)
			}

			// Calculate item total
			itemTotal := product.Price * float64(item.Quantity)
			total += itemTotal

			// Add to order items
			orderItems = append(orderItems, OrderItem{
				ProductID: productID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})

			// Update stock level
			_, err = handlers2.productsCollection.UpdateOne(
				sessionContext,
				bson.M{"_id": productID},
				bson.M{"$inc": bson.M{"stock_level": -item.Quantity}, "$set": bson.M{"updated_at": time.Now()}},
			)
			if err != nil {
				return err
			}
		}

		// Create order
		now := time.Now()
		order = Order{
			ID:        primitive.NewObjectID(),
			UserID:    req.UserID,
			Status:    "pending",
			Total:     total,
			Items:     orderItems,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Insert order
		_, err = handlers2.ordersCollection.InsertOne(sessionContext, order)
		if err != nil {
			return err
		}

		// Commit the transaction
		return session.CommitTransaction(sessionContext)
	})

	if err != nil {
		// Transaction failed, handle the error
		session.AbortTransaction(ctx)
		handlers2.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	handlers2.RespondWithJSON(w, http.StatusCreated, order)
}
