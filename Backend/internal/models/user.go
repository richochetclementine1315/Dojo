package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system.(The main user entity)
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey; default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"type:varchar(255)" json:"-"`
	AvatarURL    string    `gorm:"type: varchar(500)" json:"avatar_url"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships with other models
	Profile       *UserProfile   `gorm:"foreignKey:UserID; constraint: OnDelete:CASCADE" json:"profile,omitempty"`
	AuthAccounts  []AuthAccount  `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
	Notes         []UserNote     `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
	Sheets        []ProblemSheet `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
	Notifications []Notification `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate is a GORM hook that is triggered to generate a UUID before creating a new User record.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the User model.
func (User) TableName() string {
	return "users"
}
