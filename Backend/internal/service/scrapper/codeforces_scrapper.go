package scrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	// Fetch user info from Codeforces API
	url := fmt.Sprintf("https://codeforces.com/api/user.info?handles=%s", username)

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Codeforces data: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	var cfResp CodeforcesResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse Codeforces response: %w", err)
	}
	if cfResp.Status != "OK" || len(cfResp.Result) == 0 {
		return nil, fmt.Errorf("invalid Codeforces response")
	}
	user := cfResp.Result[0]
	stats := &PlatformStats{
		Rating:    user.Rating,
		MaxRating: user.MaxRating,

		ProblemsSolved: 0,
	}
	return stats, nil
}
