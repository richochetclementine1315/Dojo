package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserProblemProgress tracks which problems a user has solved
type UserProblemProgress struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	ProblemID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"problem_id"`
	IsSolved    bool       `gorm:"default:false;not null" json:"is_solved"`
	SolvedAt    *time.Time `json:"solved_at"`
	Attempts    int        `gorm:"default:0" json:"attempts"`
	LastAttempt *time.Time `json:"last_attempt"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Problem Problem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"problem,omitempty"`
}

// BeforeCreate hook
func (upp *UserProblemProgress) BeforeCreate(tx *gorm.DB) error {
	if upp.ID == uuid.Nil {
		upp.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (UserProblemProgress) TableName() string {
	return "user_problem_progress"
}
