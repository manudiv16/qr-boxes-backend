package services

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qr-boxes/backend/internal/models"
	"github.com/qr-boxes/backend/internal/repository"
	"github.com/skip2/go-qrcode"
)

type QRService struct {
	baseURL    string
	boxRepo    *repository.BoxRepository
}

func NewQRService(baseURL string, boxRepo *repository.BoxRepository) *QRService {
	return &QRService{
		baseURL: baseURL,
		boxRepo: boxRepo,
	}
}

func (s *QRService) CreateBox(userID string, request *models.CreateBoxRequest) (*models.CreateBoxResponse, error) {
	// Generate unique ID for the box
	boxID := uuid.New().String()

	// Create the QR code content (URL that will redirect to the box details)
	qrContent := fmt.Sprintf("%s/box/%s", s.baseURL, boxID)

	// Generate QR code as PNG bytes
	qrBytes, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Convert to base64 for easy transmission
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrBytes)

	// Process items list from text input
	items := processItemsList(request.Items)

	// Create the box object
	now := time.Now()
	box := &models.Box{
		ID:        boxID,
		UserID:    userID,
		Name:      request.Name,
		Items:     items,
		QRCode:    qrCodeBase64,
		QRCodeURL: qrContent,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Generate SVG version for web display
	qrSVG, err := s.generateQRCodeSVG(qrContent)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code SVG: %w", err)
	}

	// Save box to database
	err = s.boxRepo.Create(box)
	if err != nil {
		return nil, fmt.Errorf("failed to save box to database: %w", err)
	}

	response := &models.CreateBoxResponse{
		Box:       box,
		QRCodeSVG: qrSVG,
		Message:   "Box created successfully with QR code",
	}

	return response, nil
}

func (s *QRService) generateQRCodeSVG(content string) (string, error) {
	// For now, we'll return a placeholder SVG
	// In a real implementation, you might want to use a library that generates SVG QR codes
	svgTemplate := `<svg width="256" height="256" xmlns="http://www.w3.org/2000/svg">
		<rect width="100%%" height="100%%" fill="white"/>
		<text x="50%%" y="50%%" text-anchor="middle" dy=".3em" font-family="Arial" font-size="12">
			QR Code Generated
		</text>
		<text x="50%%" y="60%%" text-anchor="middle" dy=".3em" font-family="Arial" font-size="8">
			%s
		</text>
	</svg>`

	return fmt.Sprintf(svgTemplate, content), nil
}

func (s *QRService) GetBoxByID(boxID string) (*models.Box, error) {
	return s.boxRepo.GetByID(boxID)
}

func (s *QRService) GetUserBoxes(userID string) ([]*models.Box, error) {
	return s.boxRepo.GetByUserID(userID)
}

func (s *QRService) UpdateBox(userID string, boxID string, request *models.UpdateBoxRequest) (*models.Box, error) {
	// Get existing box to verify ownership
	box, err := s.boxRepo.GetByID(boxID)
	if err != nil {
		return nil, err
	}

	if box.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Update box data
	if request.Name != "" {
		box.Name = request.Name
	}

	if request.Items != "" {
		box.Items = processItemsList(request.Items)
	}

	// Save updated box
	err = s.boxRepo.Update(box)
	if err != nil {
		return nil, err
	}

	return box, nil
}

func (s *QRService) DeleteBox(userID string, boxID string) error {
	return s.boxRepo.Delete(boxID, userID)
}

func (s *QRService) GetUserBoxCount(userID string) (int, error) {
	return s.boxRepo.GetUserBoxCount(userID)
}

// AddItemToBox adds a single item to an existing box
func (s *QRService) AddItemToBox(userID string, boxID string, item string) (*models.Box, error) {
	// Get existing box to verify ownership
	box, err := s.boxRepo.GetByID(boxID)
	if err != nil {
		return nil, fmt.Errorf("box not found")
	}

	if box.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Add the new item to the existing items
	box.Items = append(box.Items, strings.TrimSpace(item))

	// Save updated box
	err = s.boxRepo.Update(box)
	if err != nil {
		return nil, err
	}

	return box, nil
}

// processItemsList converts a multi-line string into a cleaned list of items
func processItemsList(itemsText string) []string {
	if strings.TrimSpace(itemsText) == "" {
		return []string{}
	}

	lines := strings.Split(itemsText, "\n")
	var items []string

	for _, line := range lines {
		// Clean each line
		cleaned := strings.TrimSpace(line)
		
		// Only add non-empty items (no need to remove prefixes since items come from frontend array)
		if cleaned != "" {
			items = append(items, cleaned)
		}
	}

	return items
}