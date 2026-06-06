package scrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// LeetCodeProblem represents a single problem from LeetCode
type LeetCodeProblem struct {
	QuestionID         string `json:"questionId"`
	QuestionFrontendID string `json:"questionFrontendId"`
	Title              string `json:"title"`
	TitleSlug          string `json:"titleSlug"`
	Difficulty         string `json:"difficulty"`
	IsPaidOnly         bool   `json:"paidOnly"`
	TopicTags          []struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"topicTags"`
	AcRate float64 `json:"acRate"`
}

// LeetCodeProblemsResponse represents the response from LeetCode API
type LeetCodeProblemsResponse struct {
	Data struct {
		ProblemsetQuestionList struct {
			Total     int               `json:"total"`
			Questions []LeetCodeProblem `json:"questions"`
		} `json:"problemsetQuestionList"`
	} `json:"data"`
}

// FetchLeetCodeProblems fetches a list of problems from LeetCode
func FetchLeetCodeProblems(limit int, skip int) ([]LeetCodeProblem, int, error) {
	if limit <= 0 {
		limit = 50
	}

	// GraphQL query to fetch problems
	query := fmt.Sprintf(`{
		"query": "{problemsetQuestionList(categorySlug: \"\", limit: %d, skip: %d, filters: {}) {total questions {questionId questionFrontendId title titleSlug difficulty paidOnly topicTags {name slug} acRate}}}"
	}`, limit, skip)

	req, err := http.NewRequest("POST", "https://leetcode.com/graphql", strings.NewReader(query))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://leetcode.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch LeetCode problems: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("LeetCode API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response: %w", err)
	}

	var problemsResp LeetCodeProblemsResponse
	if err := json.Unmarshal(body, &problemsResp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse LeetCode response: %w", err)
	}

	return problemsResp.Data.ProblemsetQuestionList.Questions, problemsResp.Data.ProblemsetQuestionList.Total, nil
}

// CodeforcesProblemsResponse represents the response from Codeforces problemset.problems API
type CodeforcesProblemsResponse struct {
	Status string `json:"status"`
	Result struct {
		Problems []struct {
			ContestID int      `json:"contestId"`
			Index     string   `json:"index"`
			Name      string   `json:"name"`
			Type      string   `json:"type"`
			Rating    int      `json:"rating"`
			Tags      []string `json:"tags"`
		} `json:"problems"`
		ProblemStatistics []struct {
			ContestID   int    `json:"contestId"`
			Index       string `json:"index"`
			SolvedCount int    `json:"solvedCount"`
		} `json:"problemStatistics"`
	} `json:"result"`
	Comment string `json:"comment"`
}

// CodeforcesProblem represents a Codeforces problem with stats
type CodeforcesProblem struct {
	ContestID   int
	Index       string
	Name        string
	Type        string
	Rating      int
	Tags        []string
	SolvedCount int
}

// cfProblemsCache caches the full Codeforces problem list to avoid hitting rate limits
var cfProblemsCache struct {
	sync.RWMutex
	problems  []CodeforcesProblem
	fetchedAt time.Time
}

// cfCacheTTL — refresh cache every 6 hours
const cfCacheTTL = 6 * time.Hour

// FetchCodeforcesProblems fetches all problems from the Codeforces problemset API.
// Results are cached for 6 hours to stay well within the 1 req/2s rate limit.
func FetchCodeforcesProblems() ([]CodeforcesProblem, error) {
	cfProblemsCache.RLock()
	if len(cfProblemsCache.problems) > 0 && time.Since(cfProblemsCache.fetchedAt) < cfCacheTTL {
		problems := cfProblemsCache.problems
		cfProblemsCache.RUnlock()
		return problems, nil
	}
	cfProblemsCache.RUnlock()

	// Acquire write lock and re-check (double-checked locking)
	cfProblemsCache.Lock()
	defer cfProblemsCache.Unlock()

	if len(cfProblemsCache.problems) > 0 && time.Since(cfProblemsCache.fetchedAt) < cfCacheTTL {
		return cfProblemsCache.problems, nil
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", "https://codeforces.com/api/problemset.problems", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; DojoApp/1.0)")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		// Return stale cache rather than fail completely
		if len(cfProblemsCache.problems) > 0 {
			fmt.Printf("Warning: CF API unreachable, using stale cache: %v\n", err)
			return cfProblemsCache.problems, nil
		}
		return nil, fmt.Errorf("failed to fetch Codeforces problems: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if len(cfProblemsCache.problems) > 0 {
			return cfProblemsCache.problems, nil
		}
		return nil, fmt.Errorf("Codeforces API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var cfResp CodeforcesProblemsResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse Codeforces response: %w", err)
	}

	if cfResp.Status != "OK" {
		return nil, fmt.Errorf("Codeforces API error: %s", cfResp.Comment)
	}

	// Build solved-count lookup map
	statsMap := make(map[string]int, len(cfResp.Result.ProblemStatistics))
	for _, stat := range cfResp.Result.ProblemStatistics {
		key := fmt.Sprintf("%d-%s", stat.ContestID, stat.Index)
		statsMap[key] = stat.SolvedCount
	}

	problems := make([]CodeforcesProblem, 0, len(cfResp.Result.Problems))
	for _, p := range cfResp.Result.Problems {
		key := fmt.Sprintf("%d-%s", p.ContestID, p.Index)
		problems = append(problems, CodeforcesProblem{
			ContestID:   p.ContestID,
			Index:       p.Index,
			Name:        p.Name,
			Type:        p.Type,
			Rating:      p.Rating,
			Tags:        p.Tags,
			SolvedCount: statsMap[key],
		})
	}

	// Store in cache
	cfProblemsCache.problems = problems
	cfProblemsCache.fetchedAt = time.Now()

	return problems, nil
}

// CFProblemsFilter holds filter options for browsing CF problems directly
type CFProblemsFilter struct {
	Tags      []string
	MinRating int
	MaxRating int
	Search    string
	Page      int
	Limit     int
}

// FilterCodeforcesProblems applies filters + pagination to the cached CF problem list.
// Returns the filtered slice for the requested page, and the total count matching the filter.
func FilterCodeforcesProblems(filter CFProblemsFilter) ([]CodeforcesProblem, int, error) {
	all, err := FetchCodeforcesProblems()
	if err != nil {
		return nil, 0, err
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	searchLower := strings.ToLower(filter.Search)

	// Build a tag set for O(1) lookup
	tagSet := make(map[string]bool, len(filter.Tags))
	for _, t := range filter.Tags {
		tagSet[strings.ToLower(t)] = true
	}

	filtered := make([]CodeforcesProblem, 0)
	for _, p := range all {
		// Rating filter
		if filter.MinRating > 0 && p.Rating < filter.MinRating {
			continue
		}
		if filter.MaxRating > 0 && p.Rating > filter.MaxRating {
			continue
		}
		// Search filter
		if searchLower != "" && !strings.Contains(strings.ToLower(p.Name), searchLower) {
			continue
		}
		// Tag filter — problem must have ALL requested tags
		if len(tagSet) > 0 {
			matched := 0
			for _, pt := range p.Tags {
				if tagSet[strings.ToLower(pt)] {
					matched++
				}
			}
			if matched < len(tagSet) {
				continue
			}
		}
		filtered = append(filtered, p)
	}

	total := len(filtered)
	start := (filter.Page - 1) * filter.Limit
	if start >= total {
		return []CodeforcesProblem{}, total, nil
	}
	end := start + filter.Limit
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}
