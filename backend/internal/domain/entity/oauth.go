package entity

import (
	"encoding/json"
)

// OAuthUserData contains user information from OAuth providers
type OAuthUserData struct {
	ProviderID string `json:"provider_id"`
	Provider   string `json:"provider"` // "google", "github", "facebook"
	Email      string `json:"email"`
	Name       string `json:"name"`
	AvatarURL  string `json:"avatar_url,omitempty"`
}

type OauthGithubUserModel struct {
	ProviderID int64  `json:"id"`
	Provider   string `json:"provider"`
	Login      string `json:"login"`
	AvatarURL  string `json:"avatar_url"`
	Email      string `json:"email"`
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

func (u *OAuthUserData) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Ambil provider ID dari "sub" atau "id"
	if sub, ok := raw["sub"].(string); ok {
		u.ProviderID = sub
	} else if id, ok := raw["id"].(string); ok {
		u.ProviderID = id
	}

	if provider, ok := raw["provider"].(string); ok {
		u.Provider = provider
	}
	if email, ok := raw["email"].(string); ok {
		u.Email = email
	}
	if name, ok := raw["name"].(string); ok {
		u.Name = name
	}
	if avatar, ok := raw["picture"].(string); ok {
		u.AvatarURL = avatar
	} else if avatar, ok := raw["avatar_url"].(string); ok {
		u.AvatarURL = avatar
	}

	return nil
}
