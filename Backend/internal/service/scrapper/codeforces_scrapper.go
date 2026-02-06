package scrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CodeforcesResponse struct {
	Status string `json:"status"`
	Result []struct {
		Handle       string `json:"handle"`
		Rating       int    `json:"rating"`
		MaxRating    int    `json:"maxRating"`
		Rank         string `json:"rank"`
		MaxRank      string `json:"maxRank"`
		Contribution int    `json:"contribution"`
	} `json:"result"`
}

// FetchCodeforcesStats fetches coding statistics from Codeforces
func FetchCodeforcesStats(username string) (*PlatformStats, error) {
	// Extract username from URL if full URL is provided
	// Examples: "codeforces.com/profile/username" -> "username"
	//           "https://codeforces.com/profile/username" -> "username"
	if strings.Contains(username, "codeforces.com") {
		// Remove protocol if present
		username = strings.TrimPrefix(username, "https://")
		username = strings.TrimPrefix(username, "http://")
		// Remove domain
		username = strings.TrimPrefix(username, "codeforces.com/")
		// Remove profile prefix if present
		username = strings.TrimPrefix(username, "profile/")
		// Remove trailing slashes
		username = strings.TrimSuffix(username, "/")
	}

	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	// Fetch user info from Codeforces API
	url := fmt.Sprintf("https://codeforces.com/api/user.info?handles=%s", username)

	// HTTP GET request with proper headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Network error - return placeholder with verification link
		return &PlatformStats{
			Rating:         0,
			MaxRating:      0,
			ProblemsSolved: 0,
		}, fmt.Errorf("Codeforces API is currently unavailable. Verify username at: https://codeforces.com/profile/%s", username)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return &PlatformStats{
			Rating:         0,
			MaxRating:      0,
			ProblemsSolved: 0,
		}, fmt.Errorf("Codeforces API rate limit exceeded (Status %d). Try again later or verify at: https://codeforces.com/profile/%s", resp.StatusCode, username)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Codeforces user '%s' not found. Verify at: https://codeforces.com/profile/%s", username, username)
	}

	if resp.StatusCode != http.StatusOK {
		return &PlatformStats{
			Rating:         0,
			MaxRating:      0,
			ProblemsSolved: 0,
		}, fmt.Errorf("Codeforces API returned status %d. Service may be down. Verify at: https://codeforces.com/profile/%s", resp.StatusCode, username)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var cfResp CodeforcesResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return &PlatformStats{
			Rating:         0,
			MaxRating:      0,
			ProblemsSolved: 0,
		}, fmt.Errorf("Codeforces API response format changed. Verify username at: https://codeforces.com/profile/%s", username)
	}

	if cfResp.Status != "OK" || len(cfResp.Result) == 0 {
		return nil, fmt.Errorf("Codeforces user '%s' not found or API error. Check: https://codeforces.com/profile/%s", username, username)
	}

	user := cfResp.Result[0]

	// Fetch user submissions to count solved problems
	submissionsURL := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=10000", username)
	submissionsResp, err := http.Get(submissionsURL)
	if err != nil {
		// If submissions fetch fails, still return basic stats
		return &PlatformStats{
			Rating:         user.Rating,
			MaxRating:      user.MaxRating,
			ProblemsSolved: 0,
		}, nil
	}
	defer submissionsResp.Body.Close()

	submissionsBody, err := io.ReadAll(submissionsResp.Body)
	if err == nil {
		var submissionsData struct {
			Status string `json:"status"`
			Result []struct {
				Problem struct {
					Name  string   `json:"name"`
					Index string   `json:"index"`
					Tags  []string `json:"tags"`
				} `json:"problem"`
				Verdict string `json:"verdict"`
			} `json:"result"`
		}

		if json.Unmarshal(submissionsBody, &submissionsData) == nil && submissionsData.Status == "OK" {
			solvedProblems := make(map[string]bool)
			for _, submission := range submissionsData.Result {
				if submission.Verdict == "OK" {
					problemKey := submission.Problem.Name + submission.Problem.Index
					solvedProblems[problemKey] = true
				}
			}

			return &PlatformStats{
				Rating:         user.Rating,
				MaxRating:      user.MaxRating,
				ProblemsSolved: len(solvedProblems),
			}, nil
		}
	}

	stats := &PlatformStats{
		Rating:         user.Rating,
		MaxRating:      user.MaxRating,
		ProblemsSolved: 0,
	}
	return stats, nil
}
