package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ProblemResponse represents the problem data returned in API responses
type ProblemResponse struct {
	ID                uuid.UUID       `json:"id"`
	Platform          string          `json:"platform"`
	PlatformProblemID string          `json:"platform_problem_id"`
	Title             string          `json:"title"`
	Slug              string          `json:"slug"`
	Difficulty        string          `json:"difficulty"`
	Tags              []string        `json:"tags"`
	AcceptanceRate    float64         `json:"acceptance_rate"`
	ProblemURL        string          `json:"problem_url"`
	Description       string          `json:"description"`
	Constraints       string          `json:"constraints"`
	Examples          json.RawMessage `json:"examples"`
	Hints             json.RawMessage `json:"hints"`
	CreatedAt         time.Time       `json:"created_at"`
}
type ProblemFilterRequest struct {
	Platform   string   `query:"platform" validate:"omitempty,oneof=leetcode codeforces codechef gfg"`
	Difficulty string   `query:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	Tags       []string `query:"tags"`
	Search     string   `query:"search"`
	Page       int      `query:"page" validate:"omitempty,min=1"`
	Limit      int      `query:"limit" validate:"omitempty,min=1,max=100"`
}

// FetchproblemsRequest represents the request payload for fetching problems from external platforms
type FetchProblemsRequest struct {
	Platform string `json:"platform" validate:"required,oneof=leetcode codeforces codechef gfg"`
	Limit    int    `json:"limit" validate:"omitempty,min=1,max=100"`
}
type CreateNoteRequest struct {
	ProblemID  uuid.UUID `json:"problem_id" validate:"required,uuid"`
	Content    string    `json:"content" validate:"required"`
	IsFavorite bool      `json:"is_favorite"`
}

// NoteResponse represents the note data returned in API responses
type NoteResponse struct {
	ID         uuid.UUID       `json:"id"`
	Problem    ProblemResponse `json:"problem"`
	Content    string          `json:"content"`
	IsFavorite bool            `json:"is_favorite"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// CreateSheetRequest represents the request payload for creating a problem sheet
type CreateSheetRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

// UpdateSheetRequest represents the request payload for updating a problem sheet
type UpdateSheetRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=255"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

// SheetResponse represents the problem sheet data returned in API responses
type SheetResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	IsPublic    bool                   `json:"is_public"`
	Problems    []SheetProblemResponse `json:"problems"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// SheetProblemResponse represents the problem data within a problem sheet
type SheetProblemResponse struct {
	ID       uuid.UUID       `json:"id"`
	Problem  ProblemResponse `json:"problem"`
	Position int             `json:"position"`
	IsSolved bool            `json:"is_solved"`
	Notes    string          `json:"notes"`
	AddedAt  time.Time       `json:"added_at"`
}

// AddProblemToSheetRequest represents the request payload for adding a problem to a sheet
type AddProblemToSheetRequest struct {
	ProblemID uuid.UUID `json:"problem_id" validate:"required,uuid"`
	Position  int       `json:"position" validate:"omitempty,min=0"`
}

// UpdateSheetProblemRequest represents the request payload for updating a problem within a sheet
type UpdateSheetProblemRequest struct {
	IsSolved *bool  `json:"is_solved"`
	Notes    string `json:"notes"`
}

// CreateProblemRequest represents the request to create a problem
type CreateProblemRequest struct {
	Platform          string          `json:"platform" validate:"required,oneof=leetcode codeforces codechef gfg"`
	PlatformProblemID string          `json:"platform_problem_id" validate:"required"`
	Title             string          `json:"title" validate:"required,min=3,max=255"`
	Slug              string          `json:"slug" validate:"required"`
	Difficulty        string          `json:"difficulty" validate:"required,oneof=easy medium hard"`
	Tags              []string        `json:"tags"`
	AcceptanceRate    float64         `json:"acceptance_rate" validate:"omitempty,min=0,max=100"`
	ProblemURL        string          `json:"problem_url" validate:"required,url"`
	Description       string          `json:"description" validate:"required"`
	Constraints       string          `json:"constraints"`
	Examples          json.RawMessage `json:"examples"`
	Hints             json.RawMessage `json:"hints"`
}

// UpdateProblemRequest represents the request to update a problem
type UpdateProblemRequest struct {
	Title          string          `json:"title" validate:"omitempty,min=3,max=255"`
	Difficulty     string          `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	Tags           []string        `json:"tags"`
	AcceptanceRate float64         `json:"acceptance_rate" validate:"omitempty,min=0,max=100"`
	ProblemURL     string          `json:"problem_url" validate:"omitempty,url"`
	Description    string          `json:"description"`
	Constraints    string          `json:"constraints"`
	Examples       json.RawMessage `json:"examples"`
	Hints          json.RawMessage `json:"hints"`
}
