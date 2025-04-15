package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    int                `json:"user_id" bson:"user_id"`
	Status    string             `json:"status" bson:"status"`
	Total     float64            `json:"total" bson:"total"`
	Items     []OrderItem        `json:"items" bson:"items"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
