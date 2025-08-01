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
)

// getMongoURI constructs MongoDB connection URI from environment variables
func getMongoURI() string {
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	db := os.Getenv("MONGO_DB")
	replicaSet := os.Getenv("MONGO_REPLICA_SET")

	if user == "" {
		user = "appuser"
	}
	if password == "" {
		password = os.Getenv("DEFAULT_APP_PASSWORD")
		if password == "" {
			log.Fatal("MONGO_PASSWORD or DEFAULT_APP_PASSWORD environment variable must be set")
		}
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "27034"
	}
	if db == "" {
		db = "appdb"
	}
	if replicaSet == "" {
		replicaSet = "rs0"
	}

	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?replicaSet=%s", user, password, host, port, db, replicaSet)
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
		log.Printf("Error finding products: %v", err)
		return
	}
	defer cursor.Close(ctx)
	fmt.Println("All products:")
	i := 1
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			log.Printf("Error decoding product: %v", err)
			continue
		}
		prettyJSON, err := json.MarshalIndent(product, "", "  ")
		if err != nil {
			log.Printf("Error formatting product: %v", err)
			continue
		}
		fmt.Printf("%d.\n%s\n", i, string(prettyJSON))
		i++
	}
	fmt.Println("---")
}

// healthHandler provides health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "healthy",
		"service":   "go-app",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func main() {
	// Get MongoDB URI from environment variables
	uri := getMongoURI()
	
	// Connect to MongoDB with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)
	
	// Start HTTP server for health checks
	go func() {
		http.HandleFunc("/health", healthHandler)
		log.Println("Starting health check server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Health check server error: %v", err)
		}
	}()
	
	// Continuously poll and display products every 3 seconds
	for {
		printProducts(client)
		time.Sleep(3 * time.Second)
	}
}
