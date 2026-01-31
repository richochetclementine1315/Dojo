package service

import (
	"dojo/internal/dto"
	"dojo/internal/models"
	"dojo/internal/repository"
	"dojo/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SheetService struct {
	sheetRepo   *repository.SheetRepository
	problemRepo *repository.ProblemRepository
}

func NewSheetService(sheetRepo *repository.SheetRepository, problemRepo *repository.ProblemRepository) *SheetService {
	return &SheetService{
		sheetRepo:   sheetRepo,
		problemRepo: problemRepo,
	}
}

// CreateSheet creates a new problem sheet
func (s *SheetService) CreateSheet(userID string, req *dto.CreateSheetRequest) (*dto.SheetResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	sheet := &models.ProblemSheet{
		UserID:      userUUID,
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    req.IsPublic,
	}

	if err := s.sheetRepo.Create(sheet); err != nil {
		return nil, err
	}

	return s.mapSheetToResponse(sheet), nil
}

// GetSheetByID retrieves a sheet by ID
func (s *SheetService) GetSheetByID(userID, sheetID string) (*dto.SheetResponse, error) {
	sheet, err := s.sheetRepo.FindByID(sheetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrSheetNotFound
		}
		return nil, err
	}

	// Check access: owner or public sheet
	if sheet.UserID.String() != userID && !sheet.IsPublic {
		return nil, utils.ErrSheetAccessDenied
	}

	return s.mapSheetToResponseWithProblems(sheet), nil
}

// GetUserSheets retrieves all sheets for a user
func (s *SheetService) GetUserSheets(userID string) ([]dto.SheetResponse, error) {
	sheets, err := s.sheetRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.SheetResponse, len(sheets))
	for i, sheet := range sheets {
		responses[i] = *s.mapSheetToResponse(&sheet)
	}

	return responses, nil
}

// GetPublicSheets retrieves all public sheets
func (s *SheetService) GetPublicSheets() ([]dto.SheetResponse, error) {
	sheets, err := s.sheetRepo.FindPublicSheets()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.SheetResponse, len(sheets))
	for i, sheet := range sheets {
		responses[i] = *s.mapSheetToResponse(&sheet)
	}

	return responses, nil
}

// UpdateSheet updates a sheet
func (s *SheetService) UpdateSheet(userID, sheetID string, req *dto.UpdateSheetRequest) (*dto.SheetResponse, error) {
	sheet, err := s.sheetRepo.FindByID(sheetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrSheetNotFound
		}
		return nil, err
	}

	// Check ownership
	if sheet.UserID.String() != userID {
		return nil, utils.ErrSheetAccessDenied
	}

	// Update fields
	if req.Name != "" {
		sheet.Name = req.Name
	}
	if req.Description != "" {
		sheet.Description = req.Description
	}
	// Always update IsPublic since bool can't be nil
	sheet.IsPublic = req.IsPublic

	if err := s.sheetRepo.Update(sheet); err != nil {
		return nil, err
	}

	return s.mapSheetToResponse(sheet), nil
}

// DeleteSheet deletes a sheet
func (s *SheetService) DeleteSheet(userID, sheetID string) error {
	sheet, err := s.sheetRepo.FindByID(sheetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrSheetNotFound
		}
		return err
	}

	// Check ownership
	if sheet.UserID.String() != userID {
		return utils.ErrSheetAccessDenied
	}

	return s.sheetRepo.Delete(sheetID)
}

// AddProblemToSheet adds a problem to a sheet
func (s *SheetService) AddProblemToSheet(userID, sheetID string, req *dto.AddProblemToSheetRequest) (*dto.SheetProblemResponse, error) {
	// Check sheet ownership
	sheet, err := s.sheetRepo.FindByID(sheetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrSheetNotFound
		}
		return nil, err
	}

	if sheet.UserID.String() != userID {
		return nil, utils.ErrSheetAccessDenied
	}

	// Check if problem exists
	_, err = s.problemRepo.FindByID(req.ProblemID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrProblemNotFound
		}
		return nil, err
	}

	// Check if problem already in sheet
	exists, err := s.sheetRepo.ExistsProblemInSheet(sheetID, req.ProblemID.String())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.ErrProblemAlreadyInSheet
	}

	// Get next position
	position, err := s.sheetRepo.GetNextPosition(sheetID)
	if err != nil {
		return nil, err
	}

	sheetUUID, _ := uuid.Parse(sheetID)
	sheetProblem := &models.SheetProblem{
		SheetID:   sheetUUID,
		ProblemID: req.ProblemID,
		Position:  position,
		IsSolved:  false,
	}

	if err := s.sheetRepo.AddProblem(sheetProblem); err != nil {
		return nil, err
	}

	// Reload with problem data
	sheetProblem, err = s.sheetRepo.FindSheetProblem(sheetID, req.ProblemID.String())
	if err != nil {
		return nil, err
	}

	return s.mapSheetProblemToResponse(sheetProblem), nil
}

// RemoveProblemFromSheet removes a problem from a sheet
func (s *SheetService) RemoveProblemFromSheet(userID, sheetID, problemID string) error {
	sheet, err := s.sheetRepo.FindByID(sheetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrSheetNotFound
		}
		return err
	}

	if sheet.UserID.String() != userID {
		return utils.ErrSheetAccessDenied
	}

	return s.sheetRepo.RemoveProblem(sheetID, problemID)
}

// UpdateSheetProblem updates a problem's status/notes within a sheet
func (s *SheetService) UpdateSheetProblem(userID, sheetID, problemID string, req *dto.UpdateSheetProblemRequest) (*dto.SheetProblemResponse, error) {
	sheet, err := s.sheetRepo.FindByID(sheetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrSheetNotFound
		}
		return nil, err
	}

	if sheet.UserID.String() != userID {
		return nil, utils.ErrSheetAccessDenied
	}

	sheetProblem, err := s.sheetRepo.FindSheetProblem(sheetID, problemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrProblemNotFound
		}
		return nil, err
	}

	// Update fields
	if req.IsSolved != nil {
		sheetProblem.IsSolved = *req.IsSolved
	}
	if req.Notes != "" {
		sheetProblem.Notes = req.Notes
	}

	if err := s.sheetRepo.UpdateSheetProblem(sheetProblem); err != nil {
		return nil, err
	}

	return s.mapSheetProblemToResponse(sheetProblem), nil
}

// mapSheetToResponse converts ProblemSheet to SheetResponse (without problems)
func (s *SheetService) mapSheetToResponse(sheet *models.ProblemSheet) *dto.SheetResponse {
	return &dto.SheetResponse{
		ID:          sheet.ID,
		Name:        sheet.Name,
		Description: sheet.Description,
		IsPublic:    sheet.IsPublic,
		CreatedAt:   sheet.CreatedAt,
		UpdatedAt:   sheet.UpdatedAt,
		Problems:    []dto.SheetProblemResponse{},
	}
}

// mapSheetToResponseWithProblems converts ProblemSheet to SheetResponse (with problems)
func (s *SheetService) mapSheetToResponseWithProblems(sheet *models.ProblemSheet) *dto.SheetResponse {
	problems := make([]dto.SheetProblemResponse, len(sheet.SheetProblems))
	for i, sp := range sheet.SheetProblems {
		problems[i] = *s.mapSheetProblemToResponse(&sp)
	}

	return &dto.SheetResponse{
		ID:          sheet.ID,
		Name:        sheet.Name,
		Description: sheet.Description,
		IsPublic:    sheet.IsPublic,
		CreatedAt:   sheet.CreatedAt,
		UpdatedAt:   sheet.UpdatedAt,
		Problems:    problems,
	}
}

// mapSheetProblemToResponse converts SheetProblem to SheetProblemResponse
func (s *SheetService) mapSheetProblemToResponse(sp *models.SheetProblem) *dto.SheetProblemResponse {
	return &dto.SheetProblemResponse{
		ID: sp.ID,
		Problem: dto.ProblemResponse{
			ID:                sp.Problem.ID,
			Platform:          sp.Problem.Platform,
			PlatformProblemID: sp.Problem.PlatformProblemID,
			Title:             sp.Problem.Title,
			Slug:              sp.Problem.Slug,
			Difficulty:        sp.Problem.Difficulty,
			Tags:              []string(sp.Problem.Tags),
			AcceptanceRate:    sp.Problem.AcceptanceRate,
			ProblemURL:        sp.Problem.ProblemURL,
			Description:       sp.Problem.Description,
			Constraints:       sp.Problem.Constraints,
			Examples:          sp.Problem.Examples,
			Hints:             sp.Problem.Hints,
			CreatedAt:         sp.Problem.CreatedAt,
		},
		Position: sp.Position,
		IsSolved: sp.IsSolved,
		Notes:    sp.Notes,
		AddedAt:  sp.AddedAt,
	}
}
