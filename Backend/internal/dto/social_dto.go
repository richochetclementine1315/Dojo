package dto

import (
	"time"

	"github.com/google/uuid"
)

// SendFriendRequest represents the payload to send a friend request
type SendFriendRequestRequest struct {
	ReceiverID uuid.UUID `json:"receiver_id" validate:"required,uuid"`
}

// FriendRequestResponse represents the friend request data returned in API responses
type FriendRequestResponse struct {
	ID        uuid.UUID    `json:"id"`
	Sender    UserResponse `json:"sender"`
	Receiver  UserResponse `json:"receiver"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// UpdateFriendRequestRequest struct represents the payload to update a friend request status
type UpdateFriendRequestRequest struct {
	Action string `json:"action" validate:"required,oneof=accept reject"`
}

// FriendResponse requests represents the friend data returned in API responses
type FriendResponse struct {
	ID        uuid.UUID    `json:"id"`
	Friend    UserResponse `json:"friend"`
	IsOnline  bool         `json:"is_online"`
	CreatedAt time.Time    `json:"created_at"`
}

// BlockUserRequest represents the payload to block a user
type BlockUserRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
}

// Notificationresponse represents the notification data returned in API responses
type NotificationResponse struct {
	ID        uuid.UUID              `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	IsRead    bool                   `json:"is_read"`
	CreatedAt time.Time              `json:"created_at"`
}

// MarkNotificationReadRequest represents the payload to mark a notification as read
type MarkNotificationReadRequest struct {
	Notification []uuid.UUID `json:"notification_ids" validate:"required,min=1"`
}
