package routes

import (
	"net/http"

	"github.com/qr-boxes/backend/internal/handlers"
	"github.com/qr-boxes/backend/internal/middleware"
	"github.com/qr-boxes/backend/pkg/utils"
	"github.com/rs/cors"
)

type Router struct {
	userHandler   *handlers.UserHandler
	healthHandler *handlers.HealthHandler
	qrHandler     *handlers.QRHandler
}

func NewRouter(userHandler *handlers.UserHandler, healthHandler *handlers.HealthHandler, qrHandler *handlers.QRHandler) *Router {
	return &Router{
		userHandler:   userHandler,
		healthHandler: healthHandler,
		qrHandler:     qrHandler,
	}
}

func (rt *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Public endpoints
	mux.HandleFunc("/api/health", rt.healthHandler.GetHealth)
	mux.HandleFunc("/api/public/box", rt.qrHandler.GetPublicBoxDetails)

	// Protected endpoints requiring authentication
	mux.HandleFunc("/api/user/profile", middleware.AuthMiddleware(rt.userHandler.GetProfile))

	// QR and Box endpoints
	mux.HandleFunc("/api/boxes", middleware.AuthMiddleware(rt.qrHandler.CreateBox))
	mux.HandleFunc("/api/boxes/list", middleware.AuthMiddleware(rt.qrHandler.GetUserBoxes))
	mux.HandleFunc("/api/boxes/details", middleware.AuthMiddleware(rt.qrHandler.GetBoxByID))
	mux.HandleFunc("/api/boxes/qr", middleware.AuthMiddleware(rt.qrHandler.GetBoxQR))
	mux.HandleFunc("/api/boxes/update", middleware.AuthMiddleware(rt.qrHandler.UpdateBox))
	mux.HandleFunc("/api/boxes/add-item", middleware.AuthMiddleware(rt.qrHandler.AddItemToBox))
	mux.HandleFunc("/api/boxes/delete", middleware.AuthMiddleware(rt.qrHandler.DeleteBox))
	mux.HandleFunc("/api/boxes/stats", middleware.AuthMiddleware(rt.qrHandler.GetUserStats))

	// Setup CORS
	config := utils.GetConfig()
	allowedOrigins := []string{config.FrontendURL}

	// Add additional origins for different environments
	if config.FrontendURL != "https://manudev.dev" {
		allowedOrigins = append(allowedOrigins, "https://manudev.dev")
	}
	if config.FrontendURL != "http://localhost:4321" {
		allowedOrigins = append(allowedOrigins, "http://localhost:4321")
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum cache age for preflight options requests
	})

	return corsMiddleware.Handler(mux)
}
