package scrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// FetchGFGStats fetches coding statistics from GeeksforGeeks
func FetchGFGStats(username string) (*PlatformStats, error) {
	fmt.Printf("DEBUG GFG: Original input: %s\n", username)

	// Extract username from URL if full URL is provided
	// Examples: "geeksforgeeks.org/user/username" -> "username"
	//           "https://auth.geeksforgeeks.org/user/username/practice" -> "username"
	if strings.Contains(username, "geeksforgeeks.org") {
		// Remove protocol if present
		username = strings.TrimPrefix(username, "https://")
		username = strings.TrimPrefix(username, "http://")
		fmt.Printf("DEBUG GFG: After protocol removal: %s\n", username)

		// Remove auth. if present
		username = strings.TrimPrefix(username, "auth.")
		// Remove www. if present
		username = strings.TrimPrefix(username, "www.")
		fmt.Printf("DEBUG GFG: After subdomain removal: %s\n", username)

		// Remove domain
		username = strings.TrimPrefix(username, "geeksforgeeks.org/")
		fmt.Printf("DEBUG GFG: After domain removal: %s\n", username)

		// Remove user prefix if present
		username = strings.TrimPrefix(username, "user/")
		fmt.Printf("DEBUG GFG: After user/ removal: %s\n", username)

		// Remove trailing parts like /practice/
		if idx := strings.Index(username, "/"); idx != -1 {
			username = username[:idx]
		}
		fmt.Printf("DEBUG GFG: After slash split: %s\n", username)

		// Remove trailing slashes
		username = strings.TrimSuffix(username, "/")
	}

	fmt.Printf("DEBUG GFG: Final username: %s\n", username)

	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	// Try GFG API endpoint
	url := fmt.Sprintf("https://practiceapi.geeksforgeeks.org/api/vr/user-profile-stats/?handle=%s", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://www.geeksforgeeks.org")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Network error - return placeholder with verification link
		return &PlatformStats{
			ProblemsSolved: 0,
		}, fmt.Errorf("GeeksforGeeks API is currently unavailable. Verify username at: https://auth.geeksforgeeks.org/user/%s/practice", username)
	}
	defer resp.Body.Close()

	// Handle auth errors gracefully
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return &PlatformStats{
			ProblemsSolved: 0,
		}, fmt.Errorf("GeeksforGeeks API requires authentication (Status %d). Verify username at: https://auth.geeksforgeeks.org/user/%s/practice", resp.StatusCode, username)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("GeeksforGeeks user '%s' not found. Check: https://auth.geeksforgeeks.org/user/%s/practice", username, username)
	}

	if resp.StatusCode != http.StatusOK {
		return &PlatformStats{
			ProblemsSolved: 0,
		}, fmt.Errorf("GFG API returned status %d. Service may be down. Verify at: https://auth.geeksforgeeks.org/user/%s/practice", resp.StatusCode, username)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	bodyStr := string(body)

	// Try parsing as JSON first
	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err == nil {
		// Extract stats from JSON if available
		stats := &PlatformStats{
			ProblemsSolved: 0,
		}

		// Try to extract total problems solved
		if totalSolved, ok := jsonData["total_problems_solved"].(float64); ok {
			stats.ProblemsSolved = int(totalSolved)
		} else if totalSolved, ok := jsonData["totalProblemsSolved"].(float64); ok {
			stats.ProblemsSolved = int(totalSolved)
		}

		// Try to extract score/ranking
		if score, ok := jsonData["overall_coding_score"].(float64); ok {
			stats.GlobalRank = int(score)
		} else if score, ok := jsonData["score"].(float64); ok {
			stats.GlobalRank = int(score)
		}

		return stats, nil
	}

	// Fallback: Parse HTML
	// Look for total problems solved
	totalSolvedRe := regexp.MustCompile(`(?i)total.*?solved.*?(\d+)`)
	if matches := totalSolvedRe.FindStringSubmatch(bodyStr); len(matches) > 1 {
		if count, err := strconv.Atoi(matches[1]); err == nil {
			return &PlatformStats{
				ProblemsSolved: count,
			}, nil
		}
	}

	// Check if user exists
	if strings.Contains(bodyStr, "404") || strings.Contains(bodyStr, "not found") || strings.Contains(bodyStr, "User not found") {
		return nil, fmt.Errorf("user not found on GeeksforGeeks")
	}

	// Return basic stats if we can't parse
	stats := &PlatformStats{
		ProblemsSolved: 0,
	}

	return stats, nil
}
