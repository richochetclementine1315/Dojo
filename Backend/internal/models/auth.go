package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthAccount struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey; default: gen_random_uuid()" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null; index" json:"user_id"`
	Provider       string     `gorm:"type:varchar(50);not null" json:"provider"` // e.g., "google", "github"
	ProviderUserID string     `gorm:"type:varchar(225);not null; uniqueIndex" json:"provider_user_id"`
	AccessToken    string     `gorm:"type:text" json:"-"`
	RefreshToken   string     `gorm:"type:text" json:"-"`
	ExpiresAt      *time.Time `json:"expires_at"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relationship with User
	User User `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate is a GORM hook that is triggered to generate a UUID before creating a new AuthAccount record.
func (a *AuthAccount) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the AuthAccount model.
func (AuthAccount) TableName() string {
	return "auth_accounts"
}

// RefreshToken represents JWT refresh tokens issued to users for session management.
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey; default: gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null; index" json:"user_id"`
	Token     string    `gorm:"type:varchar(500);not null; uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationship with User
	User User `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate is a GORM hook that is triggered to generate a UUID before creating a new RefreshToken record.
func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the RefreshToken model.
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
