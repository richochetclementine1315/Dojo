package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Friend represents a friendship relationship between two users.
type Friend struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey; default: gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid; not null; index" json:"user_id"`
	FriendID  uuid.UUID `gorm:"type:uuid; not null; index" json:"friend_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationship with User
	User   User `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
	Friend User `gorm:"foreignKey:FriendID; constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate is a GORM hook that is triggered to generate a UUID before creating a new Friend record.
func (f *Friend) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the Friend model.
func (Friend) TableName() string {
	return "friends"
}

// FriendRequest represents a friend request sent from one user to another.
type FriendRequest struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey; default: gen_random_uuid()" json:"id"`
	SenderID   uuid.UUID `gorm:"type:uuid; not null; index" json:"sender_id"`
	ReceiverID uuid.UUID `gorm:"type:uuid; not null; index" json:"receiver_id"`
	Status     string    `gorm:"type:varchar(20); not null; default:'pending'" json:"status"` // e.g., "pending", "accepted", "rejected"
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship with User
	Sender   User `gorm:"foreignKey:SenderID; constraint:OnDelete:CASCADE" json:"-"`
	Receiver User `gorm:"foreignKey:ReceiverID; constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate is a GORM hook that is triggered to generate a UUID before creating a new FriendRequest record.
func (fr *FriendRequest) BeforeCreate(tx *gorm.DB) error {
	if fr.ID == uuid.Nil {
		fr.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the FriendRequest model.
func (FriendRequest) TableName() string {
	return "friend_requests"
}

// BlockedUser represents a blocking relationship where one user has blocked another.
type BlockedUser struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BlockerID uuid.UUID `gorm:"type:uuid;not null;index" json:"blocker_id"`
	BlockedID uuid.UUID `gorm:"type:uuid;not null;index" json:"blocked_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Blocker User `gorm:"foreignKey:BlockerID;constraint:OnDelete:CASCADE" json:"-"`
	Blocked User `gorm:"foreignKey:BlockedID;constraint:OnDelete:CASCADE" json:"blocked,omitempty"`
}

// BeforeCreate hook
func (b *BlockedUser) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (BlockedUser) TableName() string {
	return "blocked_users"
}

// Notification represents a notification sent to a user.
type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey; default: gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid; not null; index" json:"user_id"`
	Type      string    `gorm:"type:varchar(50); not null" json:"type"` // e.g., "friend_request", "message", etc.
	Title     string    `gorm:"type:varchar(255); not null" json:"title"`
	Message   string    `gorm:"type:text; not null" json:"message"`
	Data      string    `gorm:"type:jsonb" json:"data"` // Additional data in JSON format
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationship with User
	User User `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate is a GORM hook that is triggered to generate a UUID before creating a new Notification record.
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the Notification model.
func (Notification) TableName() string {
	return "notifications"
}
