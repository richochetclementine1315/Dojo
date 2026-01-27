package scrapper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FetchGFGStats fetches coding statistics from GeeksforGeeks
func FetchGFGStats(username string) (*PlatformStats, error) {
	// GFG profile URL
	url := fmt.Sprintf("https://auth.geeksforgeeks.org/user/%s", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GFG data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GFG returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	bodyStr := string(body)

	// Basic check if user exists
	if strings.Contains(bodyStr, "404") || strings.Contains(bodyStr, "not found") {
		return nil, fmt.Errorf("user not found on GeeksforGeeks")
	}

	stats := &PlatformStats{
		ProblemsSolved:     0,
		EasyProblemsSolved: 0,
		MedProblemsSolved:  0,
		HardProblemsSolved: 0,
		// Note: Would need HTML parsing to extract actual stats
	}

	return stats, nil
}
