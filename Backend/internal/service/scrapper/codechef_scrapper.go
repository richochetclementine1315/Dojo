package scrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CodeChefAPIResponse represents CodeChef API response
type CodeChefAPIResponse struct {
	Success bool `json:"success"`
	Result  struct {
		Data struct {
			Content struct {
				Rating        int    `json:"rating"`
				HighestRating int    `json:"highest_rating"`
				GlobalRank    int    `json:"global_rank"`
				CountryRank   int    `json:"country_rank"`
				Stars         string `json:"stars"`
			} `json:"content"`
		} `json:"data"`
	} `json:"result"`
}

// FetchCodeChefStats fetches coding statistics from CodeChef
func FetchCodeChefStats(username string) (*PlatformStats, error) {
	// Extract username from URL if full URL is provided
	// Examples: "codechef.com/users/username" -> "username"
	//           "https://www.codechef.com/users/username" -> "username"
	if strings.Contains(username, "codechef.com") {
		// Remove protocol if present
		username = strings.TrimPrefix(username, "https://")
		username = strings.TrimPrefix(username, "http://")
		// Remove www. if present
		username = strings.TrimPrefix(username, "www.")
		// Remove domain
		username = strings.TrimPrefix(username, "codechef.com/")
		// Remove users prefix if present
		username = strings.TrimPrefix(username, "users/")
		// Remove trailing slashes
		username = strings.TrimSuffix(username, "/")
	}

	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	// Try CodeChef API first (if available)
	url := fmt.Sprintf("https://codechef-api.vercel.app/handle/%s", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Network error - return placeholder with helpful message
		return &PlatformStats{
			Rating:         0,
			MaxRating:      0,
			GlobalRank:     0,
			ProblemsSolved: 0,
		}, fmt.Errorf("CodeChef API is currently unavailable. Please verify username '%s' manually at https://www.codechef.com/users/%s", username, username)
	}
	defer resp.Body.Close()

	// If API is unavailable or returns error, return basic stats instead of failing
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return &PlatformStats{
			Rating:         0,
			MaxRating:      0,
			GlobalRank:     0,
			ProblemsSolved: 0,
		}, fmt.Errorf("CodeChef API requires authentication (Status %d). Verify username at: https://www.codechef.com/users/%s", resp.StatusCode, username)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("CodeChef user '%s' not found. Check username at: https://www.codechef.com/users/%s", username, username)
	}

	if resp.StatusCode != http.StatusOK {
		return &PlatformStats{
			Rating:         0,
			MaxRating:      0,
			GlobalRank:     0,
			ProblemsSolved: 0,
		}, fmt.Errorf("CodeChef API returned status %d. Service may be down. Verify at: https://www.codechef.com/users/%s", resp.StatusCode, username)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp CodeChefAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		// If API parsing fails, return basic stats with helpful message
		return &PlatformStats{
			Rating:         0,
			MaxRating:      0,
			ProblemsSolved: 0,
		}, fmt.Errorf("CodeChef API response format changed. Verify username at: https://www.codechef.com/users/%s", username)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("CodeChef user '%s' not found", username)
	}

	stats := &PlatformStats{
		Rating:         apiResp.Result.Data.Content.Rating,
		MaxRating:      apiResp.Result.Data.Content.HighestRating,
		GlobalRank:     apiResp.Result.Data.Content.GlobalRank,
		ProblemsSolved: 0, // CodeChef API doesn't provide this directly
	}

	return stats, nil
}
