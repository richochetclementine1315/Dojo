package service

import (
	"dojo/internal/dto"
	"dojo/internal/models"
	"dojo/internal/repository"
	"dojo/internal/service/scrapper"
	"dojo/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetProfile retrieves the full user profile with platform stats
func (s *UserService) GetProfile(userID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	// Load relationships
	if err := s.userRepo.LoadProfile(user); err != nil {
		return nil, err
	}
	if err := s.userRepo.LoadPlatformStats(user); err != nil {
		return nil, err
	}

	return s.mapUserToResponse(user), nil
}

// UpdateProfile updates the user's profile information
func (s *UserService) UpdateProfile(userID string, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	// Load profile
	if err := s.userRepo.LoadProfile(user); err != nil {
		return nil, err
	}

	// Update profile fields
	if user.Profile == nil {
		return nil, errors.New("user profile not found")
	}

	profile := user.Profile
	profile.Bio = req.Bio
	profile.Location = req.Location
	profile.Website = req.Website
	profile.LeetcodeUsername = req.LeetcodeUsername
	profile.CodeforcesUsername = req.CodeforcesUsername
	profile.CodechefUsername = req.CodechefUsername
	profile.GFGUsername = req.GFGUsername

	// Update in database
	if err := s.userRepo.UpdateProfile(profile); err != nil {
		return nil, err
	}

	// Reload user with updated profile
	user, err = s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if err := s.userRepo.LoadProfile(user); err != nil {
		return nil, err
	}

	return s.mapUserToResponse(user), nil
}

// UpdateUser updates the user's account details (username, avatar)
func (s *UserService) UpdateUser(userID string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	// Check if username is already taken (if being changed)
	if req.Username != "" && req.Username != user.Username {
		existing, _ := s.userRepo.FindByUsername(req.Username)
		if existing != nil {
			return nil, utils.ErrUsernameTaken
		}
		user.Username = req.Username
	}

	// Update avatar URL if provided
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	// Update in database
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Reload user with profile
	user, err = s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if err := s.userRepo.LoadProfile(user); err != nil {
		return nil, err
	}

	return s.mapUserToResponse(user), nil
}

// ChangePassword changes the user's password
func (s *UserService) ChangePassword(userID string, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrUserNotFound
		}
		return err
	}

	// Check if user has a password (OAuth users might not)
	if user.PasswordHash == "" {
		return errors.New("cannot change password for OAuth-only accounts")
	}

	// Verify old password
	if !utils.ComparePassword(user.PasswordHash, req.OldPassword) {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}

// SyncPlatformStats syncs platform statistics for the user from external platforms
func (s *UserService) SyncPlatformStats(userID string, platforms []string) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	// Load profile to get platform usernames
	if err := s.userRepo.LoadProfile(user); err != nil {
		return nil, err
	}

	if user.Profile == nil {
		return nil, errors.New("user profile not found")
	}

	results := make(map[string]interface{})

	for _, platform := range platforms {
		switch platform {
		case "leetcode":
			if user.Profile.LeetcodeUsername != "" {
				stats, err := scrapper.FetchLeetCodeStats(user.Profile.LeetcodeUsername)
				if err != nil {
					results[platform] = map[string]string{"error": err.Error()}
					continue
				}
				// Save or update stats in database
				if err := s.savePlatformStats(userID, platform, stats); err != nil {
					results[platform] = map[string]string{"error": err.Error()}
					continue
				}
				results[platform] = map[string]string{"status": "success"}
			} else {
				results[platform] = map[string]string{"error": "LeetCode username not set"}
			}

		case "codeforces":
			if user.Profile.CodeforcesUsername != "" {
				stats, err := scrapper.FetchCodeforcesStats(user.Profile.CodeforcesUsername)
				if err != nil {
					results[platform] = map[string]string{"error": err.Error()}
					continue
				}
				if err := s.savePlatformStats(userID, platform, stats); err != nil {
					results[platform] = map[string]string{"error": err.Error()}
					continue
				}
				results[platform] = map[string]string{"status": "success"}
			} else {
				results[platform] = map[string]string{"error": "Codeforces username not set"}
			}

		case "codechef":
			if user.Profile.CodechefUsername != "" {
				stats, err := scrapper.FetchCodeChefStats(user.Profile.CodechefUsername)
				if err != nil {
					results[platform] = map[string]string{"error": err.Error()}
					continue
				}
				if err := s.savePlatformStats(userID, platform, stats); err != nil {
					results[platform] = map[string]string{"error": err.Error()}
					continue
				}
				results[platform] = map[string]string{"status": "success"}
			} else {
				results[platform] = map[string]string{"error": "CodeChef username not set"}
			}

		case "gfg":
			if user.Profile.GFGUsername != "" {
				stats, err := scrapper.FetchGFGStats(user.Profile.GFGUsername)
				if err != nil {
					results[platform] = map[string]string{"error": err.Error()}
					continue
				}
				if err := s.savePlatformStats(userID, platform, stats); err != nil {
					results[platform] = map[string]string{"error": err.Error()}
					continue
				}
				results[platform] = map[string]string{"status": "success"}
			} else {
				results[platform] = map[string]string{"error": "GFG username not set"}
			}

		default:
			results[platform] = map[string]string{"error": "unsupported platform"}
		}
	}

	return results, nil
}

// savePlatformStats saves or updates platform statistics
func (s *UserService) savePlatformStats(userID string, platform string, stats *scrapper.PlatformStats) error {
	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	platformStat := &models.UserPlatformStat{
		UserID:             userUUID,
		Platform:           platform,
		Rating:             stats.Rating,
		MaxRating:          stats.MaxRating,
		ProblemsSolved:     stats.ProblemsSolved,
		EasyProblemsSolved: stats.EasyProblemsSolved,
		MedProblemsSolved:  stats.MedProblemsSolved,
		HardProblemsSolved: stats.HardProblemsSolved,
		ContestsAttended:   stats.ContestsAttended,
		GlobalRank:         stats.GlobalRank,
	}

	return s.userRepo.UpsertPlatformStat(platformStat)
}

// mapUserToResponse converts User model to UserResponse DTO
func (s *UserService) mapUserToResponse(user *models.User) *dto.UserResponse {
	response := &dto.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Username:   user.Username,
		AvatarURL:  user.AvatarURL,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
	}

	if user.Profile != nil {
		response.Profile = &dto.ProfileResponse{
			Bio:                user.Profile.Bio,
			Location:           user.Profile.Location,
			Website:            user.Profile.Website,
			LeetcodeUsername:   user.Profile.LeetcodeUsername,
			CodeforcesUsername: user.Profile.CodeforcesUsername,
			CodechefUsername:   user.Profile.CodechefUsername,
			GFGUsername:        user.Profile.GFGUsername,
			TotalSolved:        user.Profile.TotalSolved,
			EasySolved:         user.Profile.EasySolved,
			MediumSolved:       user.Profile.MediumSolved,
			HardSolved:         user.Profile.HardSolved,
		}

		// Add platform stats if available
		if len(user.PlatformStats) > 0 {
			response.Profile.PlatformStats = make([]dto.PlatformStatResponse, len(user.PlatformStats))
			for i, stat := range user.PlatformStats {
				response.Profile.PlatformStats[i] = dto.PlatformStatResponse{
					Platform:      stat.Platform,
					Rating:        stat.Rating,
					MaxRating:     stat.MaxRating,
					SolvedCount:   stat.ProblemsSolved,
					ContestRating: 0,
					GlobalRank:    stat.GlobalRank,
					LastSyncedAt:  &stat.LastSynced,
				}
			}
		}
	}

	return response
}
