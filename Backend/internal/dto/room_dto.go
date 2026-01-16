package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateRoomRequest represents the request payload for creating a new collaborative coding room
type CreateRoomRequest struct {
	Name            string `json:"name" validate:"required,min=3,max=255"`
	MaxParticipants int    `json:"max_participants" validate:"omitempty,min=2,max=10"`
}

// JoinRoomRequest represents the request payload for joining a collaborative coding room
type JoinRoomRequest struct {
	RoomCode string `json:"room_code" validate:"required"`
}

// RoomResponse represents the room data returned in API responses
type RoomResponse struct {
	ID              uuid.UUID             `json:"id"`
	Name            string                `json:"name"`
	RoomCode        string                `json:"room_code"`
	Creator         *UserResponse         `json:"creator"`
	MaxParticipants int                   `json:"max_participants"`
	IsActive        bool                  `json:"is_active"`
	Participants    []ParticipantResponse `json:"participants"`
	CreatedAt       time.Time             `json:"created_at"`
}

// ParticipantResponse represents a participant in a collaborative coding room
type ParticipantResponse struct {
	ID       uuid.UUID    `json:"id"`
	User     UserResponse `json:"user"`
	IsOnline bool         `json:"is_online"`
	JoinedAt time.Time    `json:"joined_at"`
}

// CodeSessionResponse represents a code session in a collaborative coding room
type CodeSessionResponse struct {
	ID        uuid.UUID        `json:"id"`
	RoomID    uuid.UUID        `json:"room_id"`
	Problem   *ProblemResponse `json:"problem,omitempty"`
	Language  string           `json:"language"`
	Code      string           `json:"code"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// UpdateCodeSessionRequest represents the request payload for updating a code session
type UpdateCodeSessionRequest struct {
	Language  string     `json:"language" validate:"omitempty,oneof=go python java cpp javascript typescript rust"`
	Code      string     `json:"code"`
	ProblemID *uuid.UUID `json:"problem_id"`
}
