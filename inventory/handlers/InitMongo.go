package handlers

import (
	"context"
	"fmt"
	"inventory/models"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client               *mongo.Client
	productsCollection   *mongo.Collection
	ordersCollection     *mongo.Collection
	categoriesCollection *mongo.Collection
)

func InitMongo(ctx context.Context) {
	// Get MongoDB connection string from environment variables or use default
	mongoURI := getEnvOrDefault("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getEnvOrDefault("MONGODB_DB", "inventory-service")

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	// Create or get database
	db := Client.Database(dbName)

	// Get collections
	productsCollection = db.Collection("products")
	categoriesCollection = db.Collection("categories")
	ordersCollection = db.Collection("orders")

	// Create indexes
	_, err = productsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(false),
	})
	if err != nil {
		log.Printf("Failed to create index on products collection: %v", err)
	}

	_, err = categoriesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Printf("Failed to create index on categories collection: %v", err)
	}

	_, err = ordersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(false),
	})
	if err != nil {
		log.Printf("Failed to create index on orders collection: %v", err)
	}

	// Insert sample data if collections are empty
	count, err := categoriesCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to count categories: %v", err)
	}

	if count == 0 {
		// Insert sample categories
		electronicsCategory := models.Category{
			ID:          primitive.NewObjectID(),
			Name:        "Electronics",
			Description: "Electronic devices and accessories",
		}

		clothingCategory := models.Category{
			ID:          primitive.NewObjectID(),
			Name:        "Clothing",
			Description: "Apparel and fashion items",
		}

		_, err = categoriesCollection.InsertOne(ctx, electronicsCategory)
		if err != nil {
			log.Printf("Failed to insert sample category: %v", err)
		}

		_, err = categoriesCollection.InsertOne(ctx, clothingCategory)
		if err != nil {
			log.Printf("Failed to insert sample category: %v", err)
		}

		// Insert sample products
		now := time.Now()
		laptop := models.Product{
			ID:          primitive.NewObjectID(),
			Name:        "Laptop",
			Description: "High-performance laptop",
			Price:       999.99,
			StockLevel:  50,
			CategoryID:  electronicsCategory.ID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		tshirt := models.Product{
			ID:          primitive.NewObjectID(),
			Name:        "T-shirt",
			Description: "Cotton t-shirt",
			Price:       19.99,
			StockLevel:  100,
			CategoryID:  clothingCategory.ID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		_, err = productsCollection.InsertOne(ctx, laptop)
		if err != nil {
			log.Printf("Failed to insert sample product: %v", err)
		}

		_, err = productsCollection.InsertOne(ctx, tshirt)
		if err != nil {
			log.Printf("Failed to insert sample product: %v", err)
		}
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
