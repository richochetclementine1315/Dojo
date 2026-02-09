package service

import (
	"dojo/internal/dto"
	"dojo/internal/models"
	"dojo/internal/repository"
	"dojo/internal/service/scrapper"
	"dojo/internal/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ProblemService struct {
	problemRepo *repository.ProblemRepository
}

func NewProblemService(problemRepo *repository.ProblemRepository) *ProblemService {
	return &ProblemService{
		problemRepo: problemRepo,
	}
}

// CreateProblem creates a new problem
func (s *ProblemService) CreateProblem(req *dto.CreateProblemRequest) (*dto.ProblemResponse, error) {
	// Chk if problem url already exixts
	exists, err := s.problemRepo.ExistsByURL(req.ProblemURL)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("problem with this URL already exists")
	}
	// Create problem models
	problem := &models.Problem{
		Platform:          req.Platform,
		PlatformProblemID: req.PlatformProblemID,
		Title:             req.Title,
		Slug:              req.Slug,
		Difficulty:        req.Difficulty,
		Tags:              pq.StringArray(req.Tags),
		AcceptanceRate:    req.AcceptanceRate,
		ProblemURL:        req.ProblemURL,
		Description:       req.Description,
		Constraints:       req.Constraints,
		Examples:          req.Examples,
		Hints:             req.Hints,
	}
	err = s.problemRepo.Create(problem)
	if err != nil {
		return nil, err
	}
	return s.mapProblemToResponse(problem), nil
}

// GetProblem retrieves a problem by ID
func (s *ProblemService) GetProblemByID(id string) (*dto.ProblemResponse, error) {
	problem, err := s.problemRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrProblemNotFound
		}
		return nil, err
	}
	return s.mapProblemToResponse(problem), nil
}

// ListProblems retrieves all the problems with filtersss and paginationssss :)
func (s *ProblemService) ListProblems(filters *dto.ProblemFilterRequest, userID string) ([]dto.ProblemResponse, int64, error) {
	// Default pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 20
	}
	// Build filter map
	filterMap := make(map[string]interface{})
	if filters.Difficulty != "" {
		filterMap["difficulty"] = filters.Difficulty
	}
	if filters.Platform != "" {
		filterMap["platform"] = filters.Platform
	}
	if filters.Search != "" {
		filterMap["search"] = filters.Search
	}
	if len(filters.Tags) > 0 {
		filterMap["tags"] = filters.Tags
	}

	problems, total, err := s.problemRepo.FindAll(filterMap, filters.Page, filters.Limit)
	if err != nil {
		return nil, 0, err
	}

	// Get user's solved problems
	var userUUID uuid.UUID
	solvedMap := make(map[uuid.UUID]bool)
	if userID != "" {
		userUUID, err = uuid.Parse(userID)
		if err == nil {
			var progressList []models.UserProblemProgress
			s.problemRepo.GetDB().Where("user_id = ? AND is_solved = ?", userUUID, true).Find(&progressList)
			for _, p := range progressList {
				solvedMap[p.ProblemID] = true
			}
		}
	}

	responses := make([]dto.ProblemResponse, len(problems))
	for i, problem := range problems {
		response := s.mapProblemToResponse(&problem)
		response.IsSolved = solvedMap[problem.ID]
		responses[i] = *response
	}

	return responses, total, nil
}

// UpdateProblem updates an existing problem
func (s *ProblemService) UpdateProblem(id string, req *dto.UpdateProblemRequest) (*dto.ProblemResponse, error) {
	problem, err := s.problemRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrProblemNotFound
		}
		return nil, err
	}

	// Update fields
	if req.Title != "" {
		problem.Title = req.Title
	}
	if req.Description != "" {
		problem.Description = req.Description
	}
	if req.Difficulty != "" {
		problem.Difficulty = req.Difficulty
	}
	if req.ProblemURL != "" {
		problem.ProblemURL = req.ProblemURL
	}
	if len(req.Tags) > 0 {
		problem.Tags = pq.StringArray(req.Tags)
	}
	if req.AcceptanceRate > 0 {
		problem.AcceptanceRate = req.AcceptanceRate
	}
	if req.Constraints != "" {
		problem.Constraints = req.Constraints
	}
	if len(req.Examples) > 0 {
		problem.Examples = req.Examples
	}
	if len(req.Hints) > 0 {
		problem.Hints = req.Hints
	}

	if err := s.problemRepo.Update(problem); err != nil {
		return nil, err
	}
	return s.mapProblemToResponse(problem), nil
}

// DeleteProblem deletes a problem by ID
func (s *ProblemService) DeleteProblem(id string) error {
	_, err := s.problemRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrProblemNotFound
		}
		return err
	}

	return s.problemRepo.Delete(id)
}

// mapProblemToResponse converts Problem model to ProblemResponse DTO
func (s *ProblemService) mapProblemToResponse(problem *models.Problem) *dto.ProblemResponse {
	return &dto.ProblemResponse{
		ID:                problem.ID,
		Platform:          problem.Platform,
		PlatformProblemID: problem.PlatformProblemID,
		Title:             problem.Title,
		Slug:              problem.Slug,
		Difficulty:        problem.Difficulty,
		Tags:              []string(problem.Tags),
		AcceptanceRate:    problem.AcceptanceRate,
		ProblemURL:        problem.ProblemURL,
		Description:       problem.Description,
		Constraints:       problem.Constraints,
		Examples:          problem.Examples,
		Hints:             problem.Hints,
		CreatedAt:         problem.CreatedAt,
	}
}

// SyncProblems imports problems from external platforms
func (s *ProblemService) SyncProblems(platform string, limit int) (int, error) {
	imported := 0

	switch strings.ToLower(platform) {
	case "leetcode":
		problems, _, err := scrapper.FetchLeetCodeProblems(limit, 0)
		if err != nil {
			return 0, fmt.Errorf("failed to fetch LeetCode problems: %w", err)
		}

		for _, p := range problems {
			// Skip paid-only problems
			if p.IsPaidOnly {
				continue
			}

			// Check if problem already exists
			exists, err := s.problemRepo.ExistsByPlatformID("leetcode", p.QuestionFrontendID)
			if err != nil {
				continue
			}
			if exists {
				continue
			}

			// Extract tags
			tags := make([]string, len(p.TopicTags))
			for i, tag := range p.TopicTags {
				tags[i] = tag.Name
			}

			// Map difficulty to lowercase
			difficulty := strings.ToLower(p.Difficulty)

			// Create problem
			problem := &models.Problem{
				Platform:          "leetcode",
				PlatformProblemID: p.QuestionFrontendID,
				Title:             p.Title,
				Slug:              p.TitleSlug,
				Difficulty:        difficulty,
				Tags:              pq.StringArray(tags),
				AcceptanceRate:    p.AcRate,
				ProblemURL:        fmt.Sprintf("https://leetcode.com/problems/%s/", p.TitleSlug),
			}

			if err := s.problemRepo.Create(problem); err != nil {
				continue
			}
			imported++
		}

	case "codeforces":
		problems, err := scrapper.FetchCodeforcesProblems()
		if err != nil {
			return 0, fmt.Errorf("failed to fetch Codeforces problems: %w", err)
		}

		// Limit the number of problems
		if limit > 0 && len(problems) > limit {
			problems = problems[:limit]
		}

		for _, p := range problems {
			platformID := fmt.Sprintf("%d%s", p.ContestID, p.Index)

			// Check if problem already exists
			exists, err := s.problemRepo.ExistsByPlatformID("codeforces", platformID)
			if err != nil {
				continue
			}
			if exists {
				continue
			}

			// Map rating to difficulty
			difficulty := "medium"
			if p.Rating > 0 {
				if p.Rating < 1200 {
					difficulty = "easy"
				} else if p.Rating >= 1900 {
					difficulty = "hard"
				}
			}

			problem := &models.Problem{
				Platform:          "codeforces",
				PlatformProblemID: platformID,
				Title:             p.Name,
				Slug:              fmt.Sprintf("%d-%s", p.ContestID, strings.ToLower(p.Index)),
				Difficulty:        difficulty,
				Tags:              pq.StringArray(p.Tags),
				AcceptanceRate:    0,
				ProblemURL:        fmt.Sprintf("https://codeforces.com/problemset/problem/%d/%s", p.ContestID, p.Index),
			}

			if err := s.problemRepo.Create(problem); err != nil {
				continue
			}
			imported++
		}

	default:
		return 0, fmt.Errorf("unsupported platform: %s", platform)
	}

	return imported, nil
}

// MarkProblemSolved marks a problem as solved for a user
func (s *ProblemService) MarkProblemSolved(userID, problemID string, isSolved bool) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	problemUUID, err := uuid.Parse(problemID)
	if err != nil {
		return fmt.Errorf("invalid problem ID: %w", err)
	}

	// Check if progress record exists
	var progress models.UserProblemProgress
	result := s.problemRepo.GetDB().Where("user_id = ? AND problem_id = ?", userUUID, problemUUID).First(&progress)

	now := time.Now()

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new progress record
			progress = models.UserProblemProgress{
				UserID:    userUUID,
				ProblemID: problemUUID,
				IsSolved:  isSolved,
				Attempts:  1,
			}
			if isSolved {
				progress.SolvedAt = &now
			}
			progress.LastAttempt = &now

			return s.problemRepo.GetDB().Create(&progress).Error
		}
		return result.Error
	}

	// Update existing progress
	progress.IsSolved = isSolved
	progress.Attempts++
	progress.LastAttempt = &now
	if isSolved && progress.SolvedAt == nil {
		progress.SolvedAt = &now
	}

	return s.problemRepo.GetDB().Save(&progress).Error
}

// GetUserSolvedCount returns the count of solved problems for a user
func (s *ProblemService) GetUserSolvedCount(userID string) (int64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %w", err)
	}

	var count int64
	err = s.problemRepo.GetDB().Model(&models.UserProblemProgress{}).
		Where("user_id = ? AND is_solved = ?", userUUID, true).
		Count(&count).Error

	return count, err
}
