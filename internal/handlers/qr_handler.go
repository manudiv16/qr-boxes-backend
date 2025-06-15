package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/qr-boxes/backend/internal/middleware"
	"github.com/qr-boxes/backend/internal/models"
	"github.com/qr-boxes/backend/internal/services"
	"github.com/qr-boxes/backend/pkg/utils"
)

type QRHandler struct {
	qrService *services.QRService
}

func NewQRHandler(qrService *services.QRService) *QRHandler {
	return &QRHandler{
		qrService: qrService,
	}
}

func (h *QRHandler) CreateBox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestError(w, "Method not allowed")
		return
	}

	// Get user ID from the authenticated context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		log.Printf("QRHandler.CreateBox: Failed to get user ID: %v", err)
		utils.UnauthorizedError(w, "Invalid authentication")
		return
	}

	// Parse request body
	var request models.CreateBoxRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("QRHandler.CreateBox: Failed to decode request: %v", err)
		utils.BadRequestError(w, "Invalid request body")
		return
	}

	// Validate request
	if request.Name == "" {
		utils.BadRequestError(w, "Box name is required")
		return
	}

	if len(request.Name) > 100 {
		utils.BadRequestError(w, "Box name must be less than 100 characters")
		return
	}

	if len(request.Items) > 1000 {
		utils.BadRequestError(w, "Items list must be less than 1000 characters")
		return
	}

	log.Printf("QRHandler.CreateBox: Creating box for user %s with name: %s", userID, request.Name)

	// Create box with QR code
	response, err := h.qrService.CreateBox(userID, &request)
	if err != nil {
		log.Printf("QRHandler.CreateBox: Failed to create box: %v", err)
		utils.InternalServerError(w, "Failed to create box")
		return
	}

	log.Printf("QRHandler.CreateBox: Successfully created box %s for user %s", response.Box.ID, userID)
	utils.CreatedResponse(w, response)
}

func (h *QRHandler) GetUserBoxes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestError(w, "Method not allowed")
		return
	}

	// Get user ID from the authenticated context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		log.Printf("QRHandler.GetUserBoxes: Failed to get user ID: %v", err)
		utils.UnauthorizedError(w, "Invalid authentication")
		return
	}

	log.Printf("QRHandler.GetUserBoxes: Fetching boxes for user %s", userID)

	// Get user's boxes
	boxes, err := h.qrService.GetUserBoxes(userID)
	if err != nil {
		log.Printf("QRHandler.GetUserBoxes: Failed to fetch boxes: %v", err)
		utils.InternalServerError(w, "Failed to fetch boxes")
		return
	}

	response := map[string]interface{}{
		"boxes": boxes,
		"count": len(boxes),
	}

	log.Printf("QRHandler.GetUserBoxes: Successfully fetched %d boxes for user %s", len(boxes), userID)
	utils.SuccessResponse(w, response)
}

func (h *QRHandler) GetBoxByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestError(w, "Method not allowed")
		return
	}

	// Extract box ID from URL path
	// For now, we'll expect it as a query parameter
	boxID := r.URL.Query().Get("id")
	if boxID == "" {
		utils.BadRequestError(w, "Box ID is required")
		return
	}

	log.Printf("QRHandler.GetBoxByID: Fetching box %s", boxID)

	// Get box by ID
	box, err := h.qrService.GetBoxByID(boxID)
	if err != nil {
		log.Printf("QRHandler.GetBoxByID: Failed to fetch box: %v", err)
		utils.NotFoundError(w, "Box not found")
		return
	}

	log.Printf("QRHandler.GetBoxByID: Successfully fetched box %s", boxID)
	utils.SuccessResponse(w, box)
}

func (h *QRHandler) UpdateBox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestError(w, "Method not allowed")
		return
	}

	// Get user ID from the authenticated context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		log.Printf("QRHandler.UpdateBox: Failed to get user ID: %v", err)
		utils.UnauthorizedError(w, "Invalid authentication")
		return
	}

	// Extract box ID from URL path
	boxID := r.URL.Query().Get("id")
	if boxID == "" {
		utils.BadRequestError(w, "Box ID is required")
		return
	}

	// Parse request body
	var request models.UpdateBoxRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("QRHandler.UpdateBox: Failed to decode request: %v", err)
		utils.BadRequestError(w, "Invalid request body")
		return
	}

	// Validate request
	if request.Name != "" && len(request.Name) > 100 {
		utils.BadRequestError(w, "Box name must be less than 100 characters")
		return
	}

	if len(request.Items) > 1000 {
		utils.BadRequestError(w, "Items list must be less than 1000 characters")
		return
	}

	log.Printf("QRHandler.UpdateBox: Updating box %s for user %s", boxID, userID)

	// Update box
	box, err := h.qrService.UpdateBox(userID, boxID, &request)
	if err != nil {
		log.Printf("QRHandler.UpdateBox: Failed to update box: %v", err)
		if err.Error() == "unauthorized" {
			utils.UnauthorizedError(w, "Box not found or unauthorized")
		} else {
			utils.InternalServerError(w, "Failed to update box")
		}
		return
	}

	log.Printf("QRHandler.UpdateBox: Successfully updated box %s for user %s", boxID, userID)
	utils.SuccessResponse(w, box)
}

func (h *QRHandler) DeleteBox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.BadRequestError(w, "Method not allowed")
		return
	}

	// Get user ID from the authenticated context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		log.Printf("QRHandler.DeleteBox: Failed to get user ID: %v", err)
		utils.UnauthorizedError(w, "Invalid authentication")
		return
	}

	// Extract box ID from URL path
	boxID := r.URL.Query().Get("id")
	if boxID == "" {
		utils.BadRequestError(w, "Box ID is required")
		return
	}

	log.Printf("QRHandler.DeleteBox: Deleting box %s for user %s", boxID, userID)

	// Delete box
	err = h.qrService.DeleteBox(userID, boxID)
	if err != nil {
		log.Printf("QRHandler.DeleteBox: Failed to delete box: %v", err)
		if err.Error() == "box not found or unauthorized" {
			utils.NotFoundError(w, "Box not found or unauthorized")
		} else {
			utils.InternalServerError(w, "Failed to delete box")
		}
		return
	}

	log.Printf("QRHandler.DeleteBox: Successfully deleted box %s for user %s", boxID, userID)
	utils.SuccessResponse(w, map[string]string{"message": "Box deleted successfully"})
}

func (h *QRHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestError(w, "Method not allowed")
		return
	}

	// Get user ID from the authenticated context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		log.Printf("QRHandler.GetUserStats: Failed to get user ID: %v", err)
		utils.UnauthorizedError(w, "Invalid authentication")
		return
	}

	log.Printf("QRHandler.GetUserStats: Fetching stats for user %s", userID)

	// Get user box count
	boxCount, err := h.qrService.GetUserBoxCount(userID)
	if err != nil {
		log.Printf("QRHandler.GetUserStats: Failed to fetch stats: %v", err)
		utils.InternalServerError(w, "Failed to fetch user statistics")
		return
	}

	stats := map[string]interface{}{
		"totalBoxes": boxCount,
		"userID":     userID,
	}

	log.Printf("QRHandler.GetUserStats: Successfully fetched stats for user %s", userID)
	utils.SuccessResponse(w, stats)
}

func (h *QRHandler) GetBoxQR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestError(w, "Method not allowed")
		return
	}

	// Get user ID from the authenticated context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		log.Printf("QRHandler.GetBoxQR: Failed to get user ID: %v", err)
		utils.UnauthorizedError(w, "Invalid authentication")
		return
	}

	// Extract box ID from URL query
	boxID := r.URL.Query().Get("id")
	if boxID == "" {
		utils.BadRequestError(w, "Box ID is required")
		return
	}

	log.Printf("QRHandler.GetBoxQR: Fetching QR for box %s, user %s", boxID, userID)

	// Get box to verify ownership and retrieve QR data
	box, err := h.qrService.GetBoxByID(boxID)
	if err != nil {
		log.Printf("QRHandler.GetBoxQR: Failed to fetch box: %v", err)
		utils.NotFoundError(w, "Box not found")
		return
	}

	// Verify ownership
	if box.UserID != userID {
		log.Printf("QRHandler.GetBoxQR: Unauthorized access attempt for box %s by user %s", boxID, userID)
		utils.UnauthorizedError(w, "Unauthorized access to box")
		return
	}

	// Return box with QR data
	log.Printf("QRHandler.GetBoxQR: Successfully retrieved QR for box %s", boxID)
	utils.SuccessResponse(w, box)
}

// GetPublicBoxDetails returns box details for public access (no authentication required)
// This is used when someone scans a QR code
func (h *QRHandler) GetPublicBoxDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestError(w, "Method not allowed")
		return
	}

	// Extract box ID from URL path
	boxID := r.URL.Query().Get("id")
	if boxID == "" {
		utils.BadRequestError(w, "Box ID is required")
		return
	}

	log.Printf("QRHandler.GetPublicBoxDetails: Fetching public box details for %s", boxID)

	// Get box by ID (public access)
	box, err := h.qrService.GetBoxByID(boxID)
	if err != nil {
		log.Printf("QRHandler.GetPublicBoxDetails: Failed to fetch box: %v", err)
		utils.NotFoundError(w, "Box not found")
		return
	}

	// Create a public response without sensitive data
	publicBox := map[string]interface{}{
		"id":        box.ID,
		"name":      box.Name,
		"items":     box.Items,
		"createdAt": box.CreatedAt,
	}

	log.Printf("QRHandler.GetPublicBoxDetails: Successfully fetched public box details for %s", boxID)
	utils.SuccessResponse(w, publicBox)
}

// AddItemToBox adds a single item to an existing box
func (h *QRHandler) AddItemToBox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestError(w, "Method not allowed")
		return
	}

	// Get user ID from the authenticated context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		log.Printf("QRHandler.AddItemToBox: Failed to get user ID: %v", err)
		utils.UnauthorizedError(w, "Invalid authentication")
		return
	}

	// Parse request body
	var request struct {
		BoxID string `json:"boxId" validate:"required"`
		Item  string `json:"item" validate:"required,min=1,max=100"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("QRHandler.AddItemToBox: Failed to decode request: %v", err)
		utils.BadRequestError(w, "Invalid request body")
		return
	}

	// Validate request
	if request.BoxID == "" {
		utils.BadRequestError(w, "Box ID is required")
		return
	}

	if request.Item == "" {
		utils.BadRequestError(w, "Item is required")
		return
	}

	if len(request.Item) > 100 {
		utils.BadRequestError(w, "Item must be less than 100 characters")
		return
	}

	log.Printf("QRHandler.AddItemToBox: Adding item to box %s for user %s", request.BoxID, userID)

	// Add item to box
	box, err := h.qrService.AddItemToBox(userID, request.BoxID, request.Item)
	if err != nil {
		log.Printf("QRHandler.AddItemToBox: Failed to add item: %v", err)
		if err.Error() == "unauthorized" || err.Error() == "box not found" {
			utils.UnauthorizedError(w, "Box not found or unauthorized")
		} else {
			utils.InternalServerError(w, "Failed to add item to box")
		}
		return
	}

	log.Printf("QRHandler.AddItemToBox: Successfully added item to box %s for user %s", request.BoxID, userID)
	utils.SuccessResponse(w, box)
}