package models

import (
	"time"
)

// Box represents a physical box with QR code
type Box struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Items     []string  `json:"items,omitempty"`
	QRCode    string    `json:"qrCode"`
	QRCodeURL string    `json:"qrCodeUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateBoxRequest represents the request to create a new box
type CreateBoxRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=100"`
	Items string `json:"items,omitempty" validate:"max=1000"`
}

// CreateBoxResponse represents the response when creating a box
type CreateBoxResponse struct {
	Box       *Box   `json:"box"`
	QRCodeSVG string `json:"qrCodeSvg"`
	Message   string `json:"message"`
}

// UpdateBoxRequest represents the request to update an existing box
type UpdateBoxRequest struct {
	Name  string `json:"name,omitempty" validate:"max=100"`
	Items string `json:"items,omitempty" validate:"max=1000"`
}