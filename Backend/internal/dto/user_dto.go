package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID         uuid.UUID        `json:"id"`
	Email      string           `json:"email"`
	Username   string           `json:"username"`
	AvatarURL  string           `json:"avatar_url"`
	IsVerified bool             `json:"is_verified"`
	CreatedAt  time.Time        `json:"created_at"`
	Profile    *ProfileResponse `json:"profile,omitempty"`
}

// ProfileResponse represents the user's profile data
type ProfileResponse struct {
	Bio                string                 `json:"bio"`
	Location           string                 `json:"location"`
	Website            string                 `json:"website"`
	LeetcodeUsername   string                 `json:"leetcode_username"`
	CodeforcesUsername string                 `json:"codeforces_username"`
	CodechefUsername   string                 `json:"codechef_username"`
	GFGUsername        string                 `json:"gfg_username"`
	TotalSolved        int                    `json:"total_solved"`
	EasySolved         int                    `json:"easy_solved"`
	MediumSolved       int                    `json:"medium_solved"`
	HardSolved         int                    `json:"hard_solved"`
	PlatformStats      []PlatformStatResponse `json:"platform_stats,omitempty"`
}

// PlatformStatResponse represents the user's statistics on a coding platform
type PlatformStatResponse struct {
	Platform      string     `json:"platform"`
	Rating        int        `json:"rating"`
	MaxRating     int        `json:"max_rating"`
	SolvedCount   int        `json:"solved_count"`
	ContestRating int        `json:"contest_rating"`
	GlobalRank    int        `json:"global_rank"`
	LastSyncedAt  *time.Time `json:"last_synced_at"`
}

// UpdateProfileRequest represents the request payload for updating user profile
type UpdateProfileRequest struct {
	Bio                string `json:"bio"`
	Location           string `json:"location"`
	Website            string `json:"website"`
	LeetcodeUsername   string `json:"leetcode_username"`
	CodeforcesUsername string `json:"codeforces_username"`
	CodechefUsername   string `json:"codechef_username"`
	GFGUsername        string `json:"gfg_username"`
}

// UpdateUserRequest represents the request payload for updating user details
type UpdateUserRequest struct {
	Username  string `json:"username" validate:"omitempty,min=3,max=50"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

// ChangePasswordRequest represents the request payload for changing user password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
