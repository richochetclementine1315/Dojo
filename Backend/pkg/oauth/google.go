package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUser represents the user information retrieved from Google OAuth
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// GoogleOAuthConfig holds Google OAuth2 configuration settings
func GoogleOAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GetGoogleUserInfo fetches user information from Google using the provided OAuth2 token
func GetGoogleUserInfo(accessToken string) (*GoogleUser, error) {
	//
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from google: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed o read response: %w", err)
	}
	var user GoogleUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}
	return &user, nil
}

// ExchangeGoogleCode exchanges the authorization code for an access token
func ExchangeGoogleCode(ctx context.Context, config *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	return token, nil
}
