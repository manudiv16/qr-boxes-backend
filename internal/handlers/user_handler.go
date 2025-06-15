package handlers

import (
	"log"
	"net/http"

	"github.com/qr-boxes/backend/internal/middleware"
	"github.com/qr-boxes/backend/internal/services"
	"github.com/qr-boxes/backend/pkg/utils"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from the authenticated context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		log.Printf("UserHandler.GetProfile: Failed to get user ID: %v", err)
		utils.UnauthorizedError(w, "Invalid authentication")
		return
	}

	log.Printf("UserHandler.GetProfile: Fetching profile for user ID: %s", userID)

	// Get user profile from service
	profile, err := h.userService.GetUserProfile(r.Context(), userID)
	if err != nil {
		log.Printf("UserHandler.GetProfile: Failed to fetch user profile: %v", err)
		utils.InternalServerError(w, "Failed to fetch user details")
		return
	}

	log.Printf("UserHandler.GetProfile: Successfully fetched profile for user: %s", userID)
	utils.SuccessResponse(w, profile)
}