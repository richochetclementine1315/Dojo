package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Contest represents a coding contest from various platforms
type Contest struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Platform          string    `gorm:"type:varchar(50);not null;index" json:"platform"`
	PlatformContestID string    `gorm:"type:varchar(255)" json:"platform_contest_id"`
	Name              string    `gorm:"type:varchar(500);not null" json:"name"`
	StartTime         time.Time `gorm:"not null;index" json:"start_time"`
	DurationSeconds   int       `gorm:"not null" json:"duration_seconds"`
	ContestURL        string    `gorm:"type:text;index" json:"contest_url"`
	Description       string    `gorm:"type:text" json:"description"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Reminders []ContestReminder `gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate hook
func (c *Contest) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (Contest) TableName() string {
	return "contests"
}

// ContestReminder represents user reminders for contests
type ContestReminder struct {
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID              uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ContestID           uuid.UUID `gorm:"type:uuid;not null;index" json:"contest_id"`
	RemindBeforeMinutes int       `gorm:"default:30" json:"remind_before_minutes"` // Notify 30 mins before
	IsNotified          bool      `gorm:"default:false" json:"is_notified"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Contest Contest `gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE" json:"contest,omitempty"`
}

// BeforeCreate hook
func (r *ContestReminder) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (ContestReminder) TableName() string {
	return "contest_reminders"
}
