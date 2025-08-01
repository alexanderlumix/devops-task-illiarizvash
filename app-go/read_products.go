package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var logger *zap.Logger

// getMongoURI constructs MongoDB connection URI from environment variables
func getMongoURI() string {
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	db := os.Getenv("MONGO_DB")
	replicaSet := os.Getenv("MONGO_REPLICA_SET")

	if host == "" {
		host = "mongo-1"
	}
	if port == "" {
		port = "27017"
	}
	if db == "" {
		db = "appdb"
	}
	if replicaSet == "" {
		replicaSet = "rs0"
	}

	// Connect without authentication for development
	uri := fmt.Sprintf("mongodb://%s:%s/%s?replicaSet=%s", host, port, db, replicaSet)
	logger.Info("MongoDB URI constructed", 
		zap.String("host", host),
		zap.String("port", port),
		zap.String("database", db),
		zap.String("replicaSet", replicaSet))
	
	return uri
}

// Product represents a product document in MongoDB
type Product struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string            `bson:"name" json:"name"`
	CreatedAt time.Time         `bson:"createdAt" json:"createdAt"`
}

// printProducts retrieves and displays all products from the database
func printProducts(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	coll := client.Database("appdb").Collection("products")
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		logger.Error("Error finding products", zap.Error(err))
		return
	}
	defer cursor.Close(ctx)
	
	fmt.Println("All products:")
	i := 1
	productCount := 0
	
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			logger.Error("Error decoding product", zap.Error(err))
			continue
		}
		prettyJSON, err := json.MarshalIndent(product, "", "  ")
		if err != nil {
			logger.Error("Error formatting product", zap.Error(err))
			continue
		}
		fmt.Printf("%d.\n%s\n", i, string(prettyJSON))
		i++
		productCount++
	}
	
	logger.Info("Products retrieved successfully", zap.Int("count", productCount))
	fmt.Println("---")
}

// healthHandler provides health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]string{
		"status":    "healthy",
		"service":   "go-app",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	json.NewEncoder(w).Encode(response)
	logger.Info("Health check requested", 
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()))
}

func main() {
	// Initialize structured logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()
	
	logger.Info("Starting Go application",
		zap.String("version", "1.0.0"),
		zap.String("environment", os.Getenv("GO_ENV")),
		zap.String("log_level", "info"))
	
	// Get MongoDB URI from environment variables
	uri := getMongoURI()
	
	// Connect to MongoDB with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	logger.Info("Connecting to MongoDB", zap.String("uri", uri))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer client.Disconnect(ctx)
	
	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		logger.Fatal("Failed to ping MongoDB", zap.Error(err))
	}
	logger.Info("Successfully connected to MongoDB")
	
	// Start HTTP server for health checks
	go func() {
		http.HandleFunc("/health", healthHandler)
		logger.Info("Starting health check server", zap.String("port", "8080"))
		if err := http.ListenAndServe(":8080", nil); err != nil {
			logger.Error("Health check server error", zap.Error(err))
		}
	}()
	
	// Continuously poll and display products every 3 seconds
	logger.Info("Starting product polling loop")
	for {
		printProducts(client)
		time.Sleep(3 * time.Second)
	}
}
