package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WhiteboardSession represents a collaborative whiteboard session
type WhiteboardSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RoomID    uuid.UUID `gorm:"type:uuid;not null;index" json:"room_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Room    Room               `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"-"`
	Strokes []WhiteboardStroke `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE" json:"strokes,omitempty"`
}

// BeforeCreate hook
func (s *WhiteboardSession) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (WhiteboardSession) TableName() string {
	return "whiteboard_sessions"
}

// WhiteboardStroke represents individual drawing strokes on the whiteboard
type WhiteboardStroke struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SessionID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"session_id"`
	UserID     *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	StrokeData string     `gorm:"type:jsonb;not null" json:"stroke_data"` // {type, points, color, width, etc.}
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Session WhiteboardSession `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE" json:"-"`
	User    *User             `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
}

// BeforeCreate hook
func (s *WhiteboardStroke) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (WhiteboardStroke) TableName() string {
	return "whiteboard_strokes"
}
