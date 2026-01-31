package scrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ContestInfo represents a contest from any platform
type ContestInfo struct {
	Name       string
	Platform   string
	StartTime  time.Time
	EndTime    time.Time
	Duration   int // in seconds
	ContestURL string
	IsVirtual  bool
	Phase      string // BEFORE, CODING, FINISHED
}

// CodeforcesContestResponse represents Codeforces contest API response
type CodeforcesContestResponse struct {
	Status string `json:"status"`
	Result []struct {
		ID                  int    `json:"id"`
		Name                string `json:"name"`
		Type                string `json:"type"`
		Phase               string `json:"phase"`
		Frozen              bool   `json:"frozen"`
		DurationSeconds     int    `json:"durationSeconds"`
		StartTimeSeconds    int64  `json:"startTimeSeconds"`
		RelativeTimeSeconds int64  `json:"relativeTimeSeconds"`
	} `json:"result"`
}

// LeetCodeContestResponse represents LeetCode contest API response
type LeetCodeContestResponse struct {
	Data struct {
		TopTwoContests []struct {
			Title     string `json:"title"`
			TitleSlug string `json:"titleSlug"`
			StartTime int64  `json:"startTime"`
			Duration  int    `json:"duration"`
		} `json:"topTwoContests"`
	} `json:"data"`
}

// FetchCodeforcesContests fetches upcoming and ongoing contests from Codeforces
func FetchCodeforcesContests() ([]ContestInfo, error) {
	url := "https://codeforces.com/api/contest.list"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Codeforces contests: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var cfResp CodeforcesContestResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse Codeforces response: %w", err)
	}

	if cfResp.Status != "OK" {
		return nil, fmt.Errorf("Codeforces API error: status %s", cfResp.Status)
	}

	var contests []ContestInfo
	now := time.Now()

	for _, contest := range cfResp.Result {
		// Only include BEFORE (upcoming) and CODING (ongoing) contests
		if contest.Phase != "BEFORE" && contest.Phase != "CODING" {
			continue
		}

		startTime := time.Unix(contest.StartTimeSeconds, 0)
		endTime := startTime.Add(time.Duration(contest.DurationSeconds) * time.Second)

		// Skip contests that have already ended
		if endTime.Before(now) {
			continue
		}

		contestInfo := ContestInfo{
			Name:       contest.Name,
			Platform:   "codeforces",
			StartTime:  startTime,
			EndTime:    endTime,
			Duration:   contest.DurationSeconds,
			ContestURL: fmt.Sprintf("https://codeforces.com/contest/%d", contest.ID),
			IsVirtual:  contest.Type == "ICPC",
			Phase:      contest.Phase,
		}

		contests = append(contests, contestInfo)
	}

	return contests, nil
}

// FetchLeetCodeContests fetches upcoming contests from LeetCode
func FetchLeetCodeContests() ([]ContestInfo, error) {
	query := `{
		"query": "{topTwoContests {title titleSlug startTime duration}}"
	}`

	req, err := http.NewRequest("POST", "https://leetcode.com/graphql", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://leetcode.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LeetCode contests: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LeetCode API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var lcResp LeetCodeContestResponse
	if err := json.Unmarshal(body, &lcResp); err != nil {
		return nil, fmt.Errorf("failed to parse LeetCode response: %w", err)
	}

	var contests []ContestInfo
	now := time.Now()

	for _, contest := range lcResp.Data.TopTwoContests {
		startTime := time.Unix(contest.StartTime, 0)
		durationSeconds := contest.Duration * 60 // LeetCode returns duration in minutes
		endTime := startTime.Add(time.Duration(durationSeconds) * time.Second)

		// Determine phase
		phase := "BEFORE"
		if now.After(startTime) && now.Before(endTime) {
			phase = "CODING"
		} else if now.After(endTime) {
			continue // Skip finished contests
		}

		contestInfo := ContestInfo{
			Name:       contest.Title,
			Platform:   "leetcode",
			StartTime:  startTime,
			EndTime:    endTime,
			Duration:   durationSeconds,
			ContestURL: fmt.Sprintf("https://leetcode.com/contest/%s", contest.TitleSlug),
			IsVirtual:  false,
			Phase:      phase,
		}

		contests = append(contests, contestInfo)
	}

	return contests, nil
}

// FetchAllContests fetches contests from all supported platforms
func FetchAllContests() ([]ContestInfo, error) {
	var allContests []ContestInfo

	// Fetch from Codeforces
	cfContests, err := FetchCodeforcesContests()
	if err != nil {
		// Log error but continue with other platforms
		fmt.Printf("Warning: Failed to fetch Codeforces contests: %v\n", err)
	} else {
		allContests = append(allContests, cfContests...)
	}

	// Fetch from LeetCode
	lcContests, err := FetchLeetCodeContests()
	if err != nil {
		// Log error but continue
		fmt.Printf("Warning: Failed to fetch LeetCode contests: %v\n", err)
	} else {
		allContests = append(allContests, lcContests...)
	}

	return allContests, nil
}
