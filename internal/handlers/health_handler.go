package handlers

import (
	"net/http"
	"time"

	"github.com/qr-boxes/backend/pkg/utils"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestError(w, "Method not allowed")
		return
	}

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
		"service":   "qr-boxes-backend",
	}

	utils.SuccessResponse(w, response)
}