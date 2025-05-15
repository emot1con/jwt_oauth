package interfaces

import (
	"auth/internal/domain/entity"
	"database/sql"
)

// OAuthUsecaseInterface defines the contract for OAuth authentication operations
type OAuthUsecaseInterface interface {
	// GoogleAuth handles OAuth authentication with Google
	GoogleAuth(code string) (*entity.JWTResponse, error)

	// GitHubAuth handles OAuth authentication with GitHub
	GitHubAuth(code string) (*entity.JWTResponse, error)

	// FacebookAuth handles OAuth authentication with Facebook
	FacebookAuth(code string) (*entity.JWTResponse, error)

	// GetOAuthUser retrieves a user authenticated via OAuth
	GetOAuthUser(tx *sql.Tx, providerID string, provider string) (*entity.User, error)

	// CreateOAuthUser creates a new user based on OAuth data
	CreateOAuthUser(tx *sql.Tx, userData *entity.OAuthUserData) (*entity.User, error)
}
