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
	// Extract username from URL if full URL is provided
	// Examples: "leetcode.com/u/username/" -> "username"
	//           "https://leetcode.com/username/" -> "username"
	if strings.Contains(username, "leetcode.com") {
		// Remove protocol if present
		username = strings.TrimPrefix(username, "https://")
		username = strings.TrimPrefix(username, "http://")
		// Remove domain
		username = strings.TrimPrefix(username, "leetcode.com/")
		// Remove /u/ prefix if present
		username = strings.TrimPrefix(username, "u/")
		// Remove trailing slashes
		username = strings.TrimSuffix(username, "/")
	}

	fmt.Printf("DEBUG SCRAPER: Fetching stats for LeetCode username: %s\n", username)

	// GraphQL query for LeetCode API
	query := fmt.Sprintf(`{
        "query": "{matchedUser(username: \"%s\") {username profile {ranking} submitStats {acSubmissionNum {difficulty count}}}}"
    }`, username)

	fmt.Printf("DEBUG SCRAPER: GraphQL Query: %s\n", query)

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

	fmt.Printf("DEBUG SCRAPER: LeetCode API Status Code: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LeetCode API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("DEBUG SCRAPER: LeetCode API Response: %s\n", string(body))

	var leetCodeResp LeetCodeResponse
	if err := json.Unmarshal(body, &leetCodeResp); err != nil {
		fmt.Printf("DEBUG SCRAPER: JSON Parse Error: %v\n", err)
		return nil, fmt.Errorf("failed to parse LeetCode response: %w", err)
	}

	fmt.Printf("DEBUG SCRAPER: Parsed MatchedUser: %+v\n", leetCodeResp.Data.MatchedUser)
	fmt.Printf("DEBUG SCRAPER: Profile Ranking: %d\n", leetCodeResp.Data.MatchedUser.Profile.Ranking)
	fmt.Printf("DEBUG SCRAPER: AcSubmissionNum count: %d\n", len(leetCodeResp.Data.MatchedUser.SubmitStats.AcSubmissionNum))

	// Parse submission stats
	stats := &PlatformStats{
		GlobalRank: leetCodeResp.Data.MatchedUser.Profile.Ranking,
	}

	for _, stat := range leetCodeResp.Data.MatchedUser.SubmitStats.AcSubmissionNum {
		fmt.Printf("DEBUG SCRAPER: Processing difficulty: %s, count: %d\n", stat.Difficulty, stat.Count)
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

	fmt.Printf("DEBUG SCRAPER: Final stats - Rank: %d, Total: %d, Easy: %d, Medium: %d, Hard: %d\n",
		stats.GlobalRank, stats.ProblemsSolved, stats.EasyProblemsSolved, stats.MedProblemsSolved, stats.HardProblemsSolved)

	return stats, nil
}
