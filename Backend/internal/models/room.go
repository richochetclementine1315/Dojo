package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Room represents a collaborative coding room
type Room struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name            string     `gorm:"type:varchar(255);not null" json:"name"`
	RoomCode        string     `gorm:"type:varchar(20);not null;uniqueIndex" json:"room_code"`
	CreatedBy       *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	MaxParticipants int        `gorm:"default:4" json:"max_participants"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Creator            *User               `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"creator,omitempty"`
	Participants       []RoomParticipant   `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"participants,omitempty"`
	CodeSessions       []CodeSession       `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"code_sessions,omitempty"`
	WhiteboardSessions []WhiteboardSession `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"whiteboard_sessions,omitempty"`
}

// BefroeCreate hook
func (r *Room) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (Room) TableName() string {
	return "rooms"
}

// RoomParticipant represents a participant in a collaborative coding room
type RoomParticipant struct {
	ID       uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RoomID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"room_id"`
	UserID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	JoinedAt time.Time  `gorm:"autoCreateTime" json:"joined_at"`
	LeftAt   *time.Time `json:"left_at"`
	IsOnline bool       `gorm:"default:true" json:"is_online"`

	// Relationships
	Room Room `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"-"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate hook
func (p *RoomParticipant) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (RoomParticipant) TableName() string {
	return "room_participants"
}

// CodeSession represents a coding session within a collaborative room
type CodeSession struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RoomID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"room_id"`
	ProblemID *uuid.UUID `gorm:"type:uuid" json:"problem_id"`
	Language  string     `gorm:"type:varchar(50);not null" json:"language"`
	Code      string     `gorm:"type:text;default:''" json:"code"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Room    Room     `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"-"`
	Problem *Problem `gorm:"foreignKey:ProblemID;constraint:OnDelete:SET NULL" json:"problem,omitempty"`
}

// BeforeCreate hook
func (c *CodeSession) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (CodeSession) TableName() string {
	return "code_sessions"
}
