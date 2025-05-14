package entity

// OAuthUserData contains user information from OAuth providers
type OAuthUserData struct {
	ProviderID string `json:"provider_id"`
	Provider   string `json:"provider"` // "google", "github", "facebook"
	Email      string `json:"email"`
	Name       string `json:"name"`
	AvatarURL  string `json:"avatar_url,omitempty"`
}

// OAuthState stores state information for OAuth flow
type OAuthState struct {
	State       string `json:"state"`
	RedirectURI string `json:"redirect_uri"`
	ExpiresAt   int64  `json:"expires_at"`
}

// OAuthConfig contains OAuth provider configuration
type OAuthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	AuthURL      string `json:"auth_url"`
	TokenURL     string `json:"token_url"`
	UserInfoURL  string `json:"user_info_url"`
	Scopes       string `json:"scopes"`
}
