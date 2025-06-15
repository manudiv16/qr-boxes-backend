package utils

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	ClerkSecretKey string
	Port           int
	FrontendURL    string
	DatabaseURL    string
}

// GlobalConfig is the application configuration
var GlobalConfig Config

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Get values from environment
	clerkKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkKey == "" {
		log.Println("Warning: CLERK_SECRET_KEY not set")
	}

	// Get port with fallback to 8080
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	// Get frontend URL with fallback
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:4321" // Default Astro development server
	}

	// Get database URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://username:password@localhost:5432/qrboxes?sslmode=disable"
	}

	GlobalConfig = Config{
		ClerkSecretKey: clerkKey,
		Port:           port,
		FrontendURL:    frontendURL,
		DatabaseURL:    databaseURL,
	}

	return GlobalConfig
}

// GetConfig returns the current configuration
func GetConfig() Config {
	return GlobalConfig
}