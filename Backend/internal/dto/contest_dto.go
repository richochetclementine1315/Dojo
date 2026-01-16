package dto

import (
	"time"

	"github.com/google/uuid"
)

// ContestResponse represents the contest data returned in API responses
type ContestResponse struct {
	ID              uuid.UUID `json:"id"`
	Platform        string    `json:"platform"`
	Name            string    `json:"name"`
	StartTime       time.Time `json:"start_time"`
	DurationSeconds int       `json:"duration_seconds"`
	ContestURL      string    `json:"contest_url"`
	Description     string    `json:"description"`
	TimeUntilStart  string    `json:"time_until_start"` // in seconds
	HasReminder     bool      `json:"has_reminder"`
}

// ContestFilterRequest represents the request payload for filtering contests
type ContestFilterRequest struct {
	Platform  string     `json:"platform" validate:"omitempty,oneof=leetcode codeforces codechef gfg"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Upcoming  bool       `json:"upcoming"` // if true, fetch only upcoming contests
}

// CreateReminderRequest represents create contest reminder request
type CreateReminderRequest struct {
	ContestID           uuid.UUID `json:"contest_id" validate:"required,uuid"`
	RemindBeforeMinutes int       `json:"remind_before_minutes" validate:"omitempty,min=1,max=1440"` // Default 30
}

// ReminderResponse represents contest reminder data
type ReminderResponse struct {
	ID                  uuid.UUID       `json:"id"`
	Contest             ContestResponse `json:"contest"`
	RemindBeforeMinutes int             `json:"remind_before_minutes"`
	IsNotified          bool            `json:"is_notified"`
	CreatedAt           time.Time       `json:"created_at"`
}
