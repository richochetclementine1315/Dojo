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

type ContestService struct {
	contestRepo *repository.ContestRepository
}

func NewContestService(contestRepo *repository.ContestRepository) *ContestService {
	return &ContestService{
		contestRepo: contestRepo,
	}
}

// ListContests retrieves contests with filters and pagination
func (s *ContestService) ListContests(filters *dto.ContestFilterRequest) ([]dto.ContestResponse, int64, error) {
	// Default pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 20
	}

	// Build filter map
	filterMap := make(map[string]interface{})
	if filters.Platform != "" {
		filterMap["platform"] = filters.Platform
	}
	if filters.Upcoming {
		filterMap["upcoming"] = true
	}
	if filters.Ongoing {
		filterMap["ongoing"] = true
	}

	contests, total, err := s.contestRepo.FindAll(filterMap, filters.Page, filters.Limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.ContestResponse, len(contests))
	for i, contest := range contests {
		responses[i] = *s.mapContestToResponse(&contest)
	}

	return responses, total, nil
}

// GetContestByID retrieves a contest by ID
func (s *ContestService) GetContestByID(id string) (*dto.ContestResponse, error) {
	contest, err := s.contestRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrContestNotFound
		}
		return nil, err
	}
	return s.mapContestToResponse(contest), nil
}

// CreateReminder creates a contest reminder for a user
func (s *ContestService) CreateReminder(userID string, req *dto.CreateReminderRequest) (*dto.ReminderResponse, error) {
	// Check if contest exists
	contest, err := s.contestRepo.FindByID(req.ContestID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrContestNotFound
		}
		return nil, err
	}

	// Check if reminder already exists
	exists, err := s.contestRepo.ExistsReminder(userID, req.ContestID.String())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("reminder already exists for this contest")
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Create reminder
	reminder := &models.ContestReminder{
		UserID:              userUUID,
		ContestID:           req.ContestID,
		RemindBeforeMinutes: req.RemindBeforeMinutes,
		IsNotified:          false,
	}

	if err := s.contestRepo.CreateReminder(reminder); err != nil {
		return nil, err
	}

	// Reload with contest data
	reminder.Contest = *contest

	return s.mapReminderToResponse(reminder), nil
}

// DeleteReminder deletes a contest reminder
func (s *ContestService) DeleteReminder(userID, reminderID string) error {
	// Check if reminder exists
	reminder, err := s.contestRepo.FindReminderByID(reminderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrReminderNotFound
		}
		return err
	}

	// Check if reminder belongs to user
	if reminder.UserID.String() != userID {
		return utils.ErrUnauthorized
	}

	return s.contestRepo.DeleteReminder(reminderID)
}

// mapContestToResponse converts Contest model to ContestResponse DTO
func (s *ContestService) mapContestToResponse(contest *models.Contest) *dto.ContestResponse {
	return &dto.ContestResponse{
		ID:              contest.ID,
		Platform:        contest.Platform,
		Name:            contest.Name,
		StartTime:       contest.StartTime,
		DurationSeconds: contest.DurationSeconds,
		ContestURL:      contest.ContestURL,
		Description:     contest.Description,
		TimeUntilStart:  "",
		HasReminder:     false,
	}
}

// mapReminderToResponse converts ContestReminder model to ReminderResponse DTO
func (s *ContestService) mapReminderToResponse(reminder *models.ContestReminder) *dto.ReminderResponse {
	return &dto.ReminderResponse{
		ID:                  reminder.ID,
		Contest:             *s.mapContestToResponse(&reminder.Contest),
		RemindBeforeMinutes: reminder.RemindBeforeMinutes,
		IsNotified:          reminder.IsNotified,
		CreatedAt:           reminder.CreatedAt,
	}
}
