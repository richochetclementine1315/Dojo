package scrapper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FetchCodeChefStats fetches coding statistics from CodeChef
func FetchCodeChefStats(username string) (*PlatformStats, error) {
	// CodeChef API endpoint (unofficial)
	url := fmt.Sprintf("https://www.codechef.com/users/%s", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CodeChef data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CodeChef returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse HTML response (CodeChef doesn't have a public API)
	// This is a simplified placeholder - actual implementation would need HTML parsing
	bodyStr := string(body)

	// Basic check if user exists
	if strings.Contains(bodyStr, "404") || strings.Contains(bodyStr, "not found") {
		return nil, fmt.Errorf("user not found on CodeChef")
	}

	stats := &PlatformStats{
		Rating:         0,
		MaxRating:      0,
		ProblemsSolved: 0,
		// Note: Would need HTML parsing to extract actual stats
	}

	return stats, nil
}
