package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MongoDB URI
	MongoDb := os.Getenv("MONGO_URI")
	if MongoDb == "" {
		log.Fatal("MONGO_URI not found in environment variables")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	// Verify connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("MongoDB ping error:", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

// Global MongoDB Client
var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	// Fetch database name from .env
	databaseName := os.Getenv("MONGO_DB_NAME")
	if databaseName == "" {
		log.Fatal("MONGO_DB_NAME not found in environment variables")
	}

	var collection *mongo.Collection = client.Database(databaseName).Collection(collectionName)
	return collection
}
