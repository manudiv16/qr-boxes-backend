package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
)

var (
	ErrNoAuthHeader     = errors.New("no authorization header")
	ErrInvalidAuthToken = errors.New("invalid auth token")
	ErrUnauthorized     = errors.New("unauthorized")
)

// SetupClerk initializes Clerk with the secret key
func SetupClerk(secretKey string) {
	if secretKey == "" {
		panic("CLERK_SECRET_KEY is not set")
	}
	clerk.SetKey(secretKey)
}

// AuthMiddleware is a middleware that validates the Clerk session token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AuthMiddleware: Processing request to %s", r.URL.Path)
		
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("AuthMiddleware: No authorization header found")
			http.Error(w, ErrNoAuthHeader.Error(), http.StatusUnauthorized)
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			truncatedHeader := authHeader
			if len(authHeader) > 20 {
				truncatedHeader = authHeader[:20]
			}
			log.Printf("AuthMiddleware: Invalid token format: %s", truncatedHeader)
			http.Error(w, ErrInvalidAuthToken.Error(), http.StatusUnauthorized)
			return
		}
		token := parts[1]
		log.Printf("AuthMiddleware: Token received, length: %d", len(token))

		// Set the Clerk secret key for the package
		secretKey := os.Getenv("CLERK_SECRET_KEY")
		if secretKey == "" {
			log.Printf("AuthMiddleware: CLERK_SECRET_KEY not found")
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}
		clerk.SetKey(secretKey)
		
		// Verify the token using Clerk's JWT verification
		log.Printf("AuthMiddleware: Verifying JWT token...")
		params := jwt.VerifyParams{
			Token: token,
		}
		
		claims, err := jwt.Verify(context.Background(), &params)
		if err != nil {
			log.Printf("AuthMiddleware: JWT verification failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("AuthMiddleware: JWT verified successfully, user ID: %s", claims.Subject)
		
		// Create a context with the user information
		ctx := context.WithValue(r.Context(), "user_claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(ctx context.Context) (string, error) {
	claims, ok := ctx.Value("user_claims").(*clerk.SessionClaims)
	if !ok || claims == nil {
		return "", ErrUnauthorized
	}
	
	return claims.Subject, nil
}

// GetUserClaims retrieves the user claims from the context
func GetUserClaims(ctx context.Context) (*clerk.SessionClaims, error) {
	claims, ok := ctx.Value("user_claims").(*clerk.SessionClaims)
	if !ok || claims == nil {
		return nil, ErrUnauthorized
	}
	return claims, nil
}