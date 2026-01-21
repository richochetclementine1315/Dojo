package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHubUser represents the user information retrieved from GitHub OAuth
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
}

// GitHubEmail represents the email information retrieved from GitHub OAuth
type GitHubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

// GitHubOAuthConfig holds the configuration for GitHub OAuth2
func GitHubOAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

// GetGitHubUserInfo retrieves user information from GitHub using the provided OAuth2 token
func GetGitHubUserInfo(accessToken string) (*GitHubUser, error) {
	// This function fetches user info from GitHub API
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	// This function adds the authorization header to the request
	req.Header.Set("Authorization", "Bearer "+accessToken)
	// This function adds the Accept header to the request
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from GitHub bro: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned non-200 status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	var user GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user info in GitHub response: %w", err)
	}
	// If email is not provided in the user info, fetch it from the emails endpoint
	if user.Email == "" {
		email, err := GetGitHubPrimaryEmail(accessToken)
		if err == nil {
			user.Email = email
		}
	}
	return &user, nil
}

// GetGitHubPrimaryEmail retrieves the primary email from the Github
func GetGitHubPrimaryEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var emails []GitHubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}
	// Find primary verified email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}
	// Fallback to first verified email if no primary found
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}
	return "", fmt.Errorf("no verified email found for GitHub user")
}

// ExchangeGitHubCode exchanges the authorization code for an access token
func ExchangeGitHubCode(ctx context.Context, config *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	return token, nil
}
