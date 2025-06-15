package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/qr-boxes/backend/internal/database"
	"github.com/qr-boxes/backend/internal/handlers"
	"github.com/qr-boxes/backend/internal/repository"
	"github.com/qr-boxes/backend/internal/routes"
	"github.com/qr-boxes/backend/internal/services"
	"github.com/qr-boxes/backend/pkg/utils"
)

func main() {
	// Load configuration
	config := utils.LoadConfig()

	// Validate required environment variables
	if config.ClerkSecretKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable must be set")
	}

	// Initialize database connection
	db, err := database.NewConnection(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database schema
	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Initialize repositories
	boxRepo := repository.NewBoxRepository(db)

	// Initialize services
	userService := services.NewUserService(config.ClerkSecretKey)
	qrService := services.NewQRService(config.FrontendURL, boxRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	healthHandler := handlers.NewHealthHandler()
	qrHandler := handlers.NewQRHandler(qrService)

	// Initialize router
	router := routes.NewRouter(userHandler, healthHandler, qrHandler)
	handler := router.SetupRoutes()

	// Configure server
	port := fmt.Sprintf(":%d", config.Port)
	server := &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	log.Printf("🚀 Server starting on http://localhost%s", port)
	log.Printf("🌐 Frontend URL: %s", config.FrontendURL)
	log.Printf("📋 API Endpoints:")
	log.Printf("   GET    /api/health         - Health check")
	log.Printf("   GET    /api/user/profile   - User profile (protected)")
	log.Printf("   POST   /api/boxes          - Create new box with QR (protected)")
	log.Printf("   GET    /api/boxes/list     - Get user's boxes (protected)")
	log.Printf("   GET    /api/boxes/details  - Get box details (protected)")
	log.Printf("   PUT    /api/boxes/update   - Update box (protected)")
	log.Printf("   DELETE /api/boxes/delete   - Delete box (protected)")
	log.Printf("   GET    /api/boxes/stats    - Get user statistics (protected)")
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}