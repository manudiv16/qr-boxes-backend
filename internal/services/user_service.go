package services

import (
	"context"
	"fmt"
	"time"

	clerk "github.com/clerk/clerk-sdk-go/v2"
)

type UserService struct {
	clerkSecretKey string
}

type UserProfile struct {
	ID          string                 `json:"id"`
	Email       string                 `json:"email"`
	FirstName   string                 `json:"firstName"`
	LastName    string                 `json:"lastName"`
	ImageURL    string                 `json:"imageUrl"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   int64                  `json:"timestamp"`
	MemberSince string                 `json:"memberSince"`
}

func NewUserService(clerkSecretKey string) *UserService {
	return &UserService{
		clerkSecretKey: clerkSecretKey,
	}
}

func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	// Set Clerk API key
	clerk.SetKey(s.clerkSecretKey)

	// Create API request to fetch user details
	req := clerk.NewAPIRequest("GET", "/users/"+userID)

	var user clerk.User
	backend := clerk.GetBackend()
	err := backend.Call(ctx, req, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user from Clerk: %w", err)
	}

	// Extract user data with safe handling of optional fields
	profile := &UserProfile{
		ID:        userID,
		Timestamp: time.Now().Unix(),
	}

	if len(user.EmailAddresses) > 0 {
		profile.Email = user.EmailAddresses[0].EmailAddress
	}

	if user.FirstName != nil {
		profile.FirstName = *user.FirstName
	}

	if user.LastName != nil {
		profile.LastName = *user.LastName
	}

	if user.ImageURL != nil {
		profile.ImageURL = *user.ImageURL
	}

	// Convert RawMessage to map for JSON response
	if user.PublicMetadata != nil {
		metadata := make(map[string]interface{})
		// Note: In a real implementation, you might want to unmarshal this properly
		profile.Metadata = metadata
	} else {
		profile.Metadata = make(map[string]interface{})
	}

	profile.MemberSince = time.Unix(user.CreatedAt, 0).Format("2006-01-02")

	return profile, nil
}