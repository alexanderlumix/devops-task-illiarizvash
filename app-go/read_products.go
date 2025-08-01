package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Rate limiter structure
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	window   time.Duration
	maxReqs  int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(window time.Duration, maxReqs int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		window:   window,
		maxReqs:  maxReqs,
	}
}

// IsAllowed checks if request is allowed
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Clean old requests
	if times, exists := rl.requests[ip]; exists {
		var validTimes []time.Time
		for _, t := range times {
			if t.After(windowStart) {
				validTimes = append(validTimes, t)
			}
		}
		rl.requests[ip] = validTimes
	}

	// Check if limit exceeded
	if len(rl.requests[ip]) >= rl.maxReqs {
		return false
	}

	// Add current request
	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

// getMongoURI constructs MongoDB connection URI from environment variables
func getMongoURI() string {
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	db := os.Getenv("MONGO_DB")
	replicaSet := os.Getenv("MONGO_REPLICA_SET")

	if host == "" {
		host = "mongo-0"  // Use primary node
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

// ProductRequest represents a product creation request
type ProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price,omitempty"`
	Description string  `json:"description,omitempty"`
}

// ValidationError represents validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// validateProduct validates product creation request
func validateProduct(req ProductRequest) []ValidationError {
	var errors []ValidationError

	// Validate name
	if strings.TrimSpace(req.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name is required",
		})
	} else if len(req.Name) > 100 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name must be less than 100 characters",
		})
	}

	// Validate price
	if req.Price < 0 {
		errors = append(errors, ValidationError{
			Field:   "price",
			Message: "Price must be non-negative",
		})
	}

	// Validate description
	if len(req.Description) > 500 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "Description must be less than 500 characters",
		})
	}

	return errors
}

// sanitizeInput sanitizes input strings
func sanitizeInput(input string) string {
	// Remove potentially dangerous characters
	input = strings.ReplaceAll(input, "<script>", "")
	input = strings.ReplaceAll(input, "</script>", "")
	input = strings.ReplaceAll(input, "javascript:", "")
	return strings.TrimSpace(input)
}

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

// Rate limiting middleware
func rateLimitMiddleware(limiter *RateLimiter) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				ip = strings.Split(forwardedFor, ",")[0]
			}

			if !limiter.IsAllowed(ip) {
				logger.Warn("Rate limit exceeded", zap.String("ip", ip))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Rate limit exceeded. Please try again later.",
				})
				return
			}

			next(w, r)
		}
	}
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

// createProductHandler handles product creation with validation
func createProductHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse request body
		var req ProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Error decoding request", zap.Error(err))
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Sanitize input
		req.Name = sanitizeInput(req.Name)
		req.Description = sanitizeInput(req.Description)

		// Validate request
		errors := validateProduct(req)
		if len(errors) > 0 {
			logger.Warn("Validation failed", 
				zap.String("remote_addr", r.RemoteAddr),
				zap.Any("errors", errors))
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":  "Validation failed",
				"errors": errors,
			})
			return
		}

		// Create product in database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		coll := client.Database("appdb").Collection("products")
		result, err := coll.InsertOne(ctx, bson.M{
			"name":      req.Name,
			"price":     req.Price,
			"description": req.Description,
			"createdAt": time.Now(),
		})

		if err != nil {
			logger.Error("Error creating product", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Product created successfully",
			"productId": result.InsertedID,
			"product": map[string]interface{}{
				"name":        req.Name,
				"price":       req.Price,
				"description": req.Description,
				"createdAt":   time.Now(),
			},
		})

		logger.Info("Product created successfully",
			zap.String("product_name", req.Name),
			zap.String("remote_addr", r.RemoteAddr))
	}
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
	
	// Initialize rate limiter (100 requests per 15 minutes)
	limiter := NewRateLimiter(15*time.Minute, 100)
	
	// Start HTTP server for health checks and product creation
	go func() {
		http.HandleFunc("/health", corsMiddleware(rateLimitMiddleware(limiter)(healthHandler)))
		http.HandleFunc("/products", corsMiddleware(rateLimitMiddleware(limiter)(createProductHandler(client))))
		
		logger.Info("Starting HTTP server", zap.String("port", "8080"))
		if err := http.ListenAndServe(":8080", nil); err != nil {
			logger.Error("HTTP server error", zap.Error(err))
		}
	}()
	
	// Continuously poll and display products every 3 seconds
	logger.Info("Starting product polling loop")
	for {
		printProducts(client)
		time.Sleep(3 * time.Second)
	}
}
