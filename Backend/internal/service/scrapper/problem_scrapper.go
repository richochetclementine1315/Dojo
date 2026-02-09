package scrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// CodeforcesProblemsResponse represents the response from Codeforces API
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
}

// CodeforcesProblem represents a Codeforces problem with stats
type CodeforcesProblem struct {
	ContestID   int
	Index       string
	Name        string
	Rating      int
	Tags        []string
	SolvedCount int
}

// FetchCodeforcesProblems fetches problems from Codeforces API
func FetchCodeforcesProblems() ([]CodeforcesProblem, error) {
	resp, err := http.Get("https://codeforces.com/api/problemset.problems")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Codeforces problems: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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
		return nil, fmt.Errorf("Codeforces API returned error status")
	}

	// Create a map of problem stats for quick lookup
	statsMap := make(map[string]int)
	for _, stat := range cfResp.Result.ProblemStatistics {
		key := fmt.Sprintf("%d-%s", stat.ContestID, stat.Index)
		statsMap[key] = stat.SolvedCount
	}

	// Combine problems with their stats
	problems := make([]CodeforcesProblem, 0, len(cfResp.Result.Problems))
	for _, p := range cfResp.Result.Problems {
		key := fmt.Sprintf("%d-%s", p.ContestID, p.Index)
		problems = append(problems, CodeforcesProblem{
			ContestID:   p.ContestID,
			Index:       p.Index,
			Name:        p.Name,
			Rating:      p.Rating,
			Tags:        p.Tags,
			SolvedCount: statsMap[key],
		})
	}

	return problems, nil
}
