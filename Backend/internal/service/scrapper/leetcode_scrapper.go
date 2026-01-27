package scrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// LeetCodeResponse represents the response from LeetCode GraphQL API
type LeetCodeResponse struct {
	Data struct {
		MatchedUser struct {
			Username string `json:"username"`
			Profile  struct {
				Ranking int `json:"ranking"`
			} `json:"profile"`
			SubmitStats struct {
				AcSubmissionNum []struct {
					Difficulty string `json:"difficulty"`
					Count      int    `json:"count"`
				} `json:"acSubmissionNum"`
			} `json:"submitStats"`
		} `json:"matchedUser"`
	} `json:"data"`
}

// FetchLeetCodeStats fetches coding statistics from LeetCode
func FetchLeetCodeStats(username string) (*PlatformStats, error) {
	// GraphQL query for LeetCode API
	query := fmt.Sprintf(`{
        "query": "{matchedUser(username: \"%s\") {username profile {ranking} submitStats {acSubmissionNum {difficulty count}}}}"
    }`, username)

	// Make HTTP request to LeetCode GraphQL API
	req, err := http.NewRequest("POST", "https://leetcode.com/graphql", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://leetcode.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LeetCode data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LeetCode API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var leetCodeResp LeetCodeResponse
	if err := json.Unmarshal(body, &leetCodeResp); err != nil {
		return nil, fmt.Errorf("failed to parse LeetCode response: %w", err)
	}

	// Parse submission stats
	stats := &PlatformStats{
		GlobalRank: leetCodeResp.Data.MatchedUser.Profile.Ranking,
	}

	for _, stat := range leetCodeResp.Data.MatchedUser.SubmitStats.AcSubmissionNum {
		stats.ProblemsSolved += stat.Count
		switch stat.Difficulty {
		case "Easy":
			stats.EasyProblemsSolved = stat.Count
		case "Medium":
			stats.MedProblemsSolved = stat.Count
		case "Hard":
			stats.HardProblemsSolved = stat.Count
		}
	}

	return stats, nil
}
