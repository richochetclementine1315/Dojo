package service

import (
	"dojo/internal/dto"
	"dojo/internal/models"
	"dojo/internal/repository"
	"dojo/internal/utils"
	"errors"

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
func (s *ProblemService) ListProblems(filters *dto.ProblemFilterRequest) ([]dto.ProblemResponse, int64, error) {
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

	responses := make([]dto.ProblemResponse, len(problems))
	for i, problem := range problems {
		responses[i] = *s.mapProblemToResponse(&problem)
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
