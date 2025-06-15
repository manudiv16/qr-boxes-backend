package utils

import (
	"encoding/json"
	"net/http"
)

// Response is a standardized API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Status  int         `json:"-"` // HTTP status code, not serialized
}

// JSONResponse sends a JSON response with appropriate headers
func JSONResponse(w http.ResponseWriter, status int, success bool, data interface{}, errMsg string) {
	response := Response{
		Success: success,
		Status:  status,
	}

	if data != nil {
		response.Data = data
	}

	if errMsg != "" {
		response.Error = errMsg
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// SuccessResponse sends a successful JSON response
func SuccessResponse(w http.ResponseWriter, data interface{}) {
	JSONResponse(w, http.StatusOK, true, data, "")
}

// CreatedResponse sends a 201 Created JSON response
func CreatedResponse(w http.ResponseWriter, data interface{}) {
	JSONResponse(w, http.StatusCreated, true, data, "")
}

// ErrorResponse sends an error JSON response
func ErrorResponse(w http.ResponseWriter, status int, errMsg string) {
	JSONResponse(w, status, false, nil, errMsg)
}

// BadRequestError sends a 400 Bad Request error response
func BadRequestError(w http.ResponseWriter, errMsg string) {
	ErrorResponse(w, http.StatusBadRequest, errMsg)
}

// UnauthorizedError sends a 401 Unauthorized error response
func UnauthorizedError(w http.ResponseWriter, errMsg string) {
	if errMsg == "" {
		errMsg = "Unauthorized"
	}
	ErrorResponse(w, http.StatusUnauthorized, errMsg)
}

// ForbiddenError sends a 403 Forbidden error response
func ForbiddenError(w http.ResponseWriter, errMsg string) {
	if errMsg == "" {
		errMsg = "Forbidden"
	}
	ErrorResponse(w, http.StatusForbidden, errMsg)
}

// NotFoundError sends a 404 Not Found error response
func NotFoundError(w http.ResponseWriter, errMsg string) {
	if errMsg == "" {
		errMsg = "Resource not found"
	}
	ErrorResponse(w, http.StatusNotFound, errMsg)
}

// ServerError sends a 500 Internal Server Error response
func ServerError(w http.ResponseWriter, errMsg string) {
	if errMsg == "" {
		errMsg = "Internal server error"
	}
	ErrorResponse(w, http.StatusInternalServerError, errMsg)
}

// InternalServerError sends a 500 Internal Server Error response (alias for ServerError)
func InternalServerError(w http.ResponseWriter, errMsg string) {
	ServerError(w, errMsg)
}