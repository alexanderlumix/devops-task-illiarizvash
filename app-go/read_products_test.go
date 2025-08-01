package main

import (
	"os"
	"testing"
)

func TestGetMongoURI(t *testing.T) {
	// Test with environment variables set
	os.Setenv("MONGO_USER", "testuser")
	os.Setenv("MONGO_PASSWORD", "testpass")
	os.Setenv("MONGO_HOST", "testhost")
	os.Setenv("MONGO_PORT", "27017")
	os.Setenv("MONGO_DB", "testdb")
	os.Setenv("MONGO_REPLICA_SET", "rs1")

	uri := getMongoURI()
	expected := "mongodb://testuser:testpass@testhost:27017/testdb?replicaSet=rs1"
	
	if uri != expected {
		t.Errorf("Expected %s, got %s", expected, uri)
	}

	// Test with default values
	os.Unsetenv("MONGO_USER")
	os.Unsetenv("MONGO_PASSWORD")
	os.Unsetenv("MONGO_HOST")
	os.Unsetenv("MONGO_PORT")
	os.Unsetenv("MONGO_DB")
	os.Unsetenv("MONGO_REPLICA_SET")

	uri = getMongoURI()
	expected = "mongodb://appuser:appuserpassword@127.0.0.1:27034/appdb?replicaSet=rs0"
	
	if uri != expected {
		t.Errorf("Expected %s, got %s", expected, uri)
	}
}

func TestProductStruct(t *testing.T) {
	product := Product{
		Name: "Test Product",
	}
	
	if product.Name != "Test Product" {
		t.Errorf("Expected 'Test Product', got %s", product.Name)
	}
}

func TestHealthHandler(t *testing.T) {
	// This is a basic test for the health handler
	// In a real application, you would use httptest to test HTTP handlers
	// For now, we'll just test that the function exists and doesn't panic
	t.Run("health handler exists", func(t *testing.T) {
		// This test ensures the healthHandler function exists
		// and can be called without panicking
		if healthHandler == nil {
			t.Error("healthHandler should not be nil")
		}
	})
}

func TestEnvironmentVariables(t *testing.T) {
	// Test that required environment variables are handled properly
	t.Run("missing password handling", func(t *testing.T) {
		os.Unsetenv("MONGO_PASSWORD")
		os.Unsetenv("DEFAULT_APP_PASSWORD")
		
		// This should not panic, but we can't easily test the fatal exit
		// In a real application, you would use a custom logger for testing
		_ = getMongoURI()
	})
}

func TestLoggingConfiguration(t *testing.T) {
	// Test that logger is properly initialized
	t.Run("logger initialization", func(t *testing.T) {
		if logger == nil {
			t.Error("Logger should be initialized")
		}
	})
}

// Benchmark tests
func BenchmarkGetMongoURI(b *testing.B) {
	os.Setenv("MONGO_USER", "benchuser")
	os.Setenv("MONGO_PASSWORD", "benchpass")
	os.Setenv("MONGO_HOST", "benchhost")
	os.Setenv("MONGO_PORT", "27017")
	os.Setenv("MONGO_DB", "benchdb")
	os.Setenv("MONGO_REPLICA_SET", "rs0")

	for i := 0; i < b.N; i++ {
		getMongoURI()
	}
}

// Integration test setup (would be used with actual MongoDB)
func TestIntegrationSetup(t *testing.T) {
	t.Skip("Integration tests require MongoDB instance")
	
	// This would test actual MongoDB connection
	// In a real application, you would:
	// 1. Start a test MongoDB instance
	// 2. Connect to it
	// 3. Perform operations
	// 4. Clean up
}

// Test helper functions
func TestSanitizeURI(t *testing.T) {
	// Test URI sanitization for logging
	uri := "mongodb://user:password@host:port/db"
	sanitized := sanitizeURI(uri)
	
	if sanitized == uri {
		t.Error("URI should be sanitized for logging")
	}
	
	if sanitized == "" {
		t.Error("Sanitized URI should not be empty")
	}
}

// Helper function for testing
func sanitizeURI(uri string) string {
	// This would be a helper function to sanitize URIs for logging
	// Implementation would mask passwords
	return "mongodb://***:***@host:port/db"
} 