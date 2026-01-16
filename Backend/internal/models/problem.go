package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Problem represents a coding problem from various platforms
type Problem struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Platform          string         `gorm:"type:varchar(50);not null;index" json:"platform"` // 'leetcode', 'codeforces', etc.
	PlatformProblemID string         `gorm:"type:varchar(255);not null" json:"platform_problem_id"`
	Title             string         `gorm:"type:varchar(500);not null" json:"title"`
	Slug              string         `gorm:"type:varchar(500)" json:"slug"`
	Difficulty        string         `gorm:"type:varchar(20);index" json:"difficulty"` // 'easy', 'medium', 'hard'
	Tags              pq.StringArray `gorm:"type:text[]" json:"tags"`                  // PostgreSQL array
	AcceptanceRate    float64        `json:"acceptance_rate"`
	ProblemURL        string         `gorm:"type:text" json:"problem_url"`
	Description       string         `gorm:"type:text" json:"description"`
	Constraints       string         `gorm:"type:text" json:"constraints"`
	Examples          string         `gorm:"type:jsonb" json:"examples"`
	Hints             string         `gorm:"type:jsonb" json:"hints"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Notes         []UserNote     `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"-"`
	SheetProblems []SheetProblem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate hook
func (p *Problem) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (Problem) TableName() string {
	return "problems"
}

// ProblemSheet represents a collection of problems created by users
type ProblemSheet struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User          User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	SheetProblems []SheetProblem `gorm:"foreignKey:SheetID;constraint:OnDelete:CASCADE" json:"problems,omitempty"`
}

// BeforeCreate hook
func (s *ProblemSheet) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (ProblemSheet) TableName() string {
	return "problem_sheets"
}

// SheetProblem represents the junction table between sheets and problems
type SheetProblem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SheetID   uuid.UUID `gorm:"type:uuid;not null;index" json:"sheet_id"`
	ProblemID uuid.UUID `gorm:"type:uuid;not null;index" json:"problem_id"`
	Position  int       `gorm:"not null" json:"position"`
	IsSolved  bool      `gorm:"default:false" json:"is_solved"`
	Notes     string    `gorm:"type:text" json:"notes"`
	AddedAt   time.Time `gorm:"autoCreateTime" json:"added_at"`

	// Relationships
	Sheet   ProblemSheet `gorm:"foreignKey:SheetID;constraint:OnDelete:CASCADE" json:"-"`
	Problem Problem      `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"problem,omitempty"`
}

// BeforeCreate hook
func (sp *SheetProblem) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (SheetProblem) TableName() string {
	return "sheet_problems"
}

// UserNote represents personal notes for problems
type UserNote struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ProblemID  uuid.UUID `gorm:"type:uuid;not null;index" json:"problem_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsFavorite bool      `gorm:"default:false" json:"is_favorite"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Problem Problem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"problem,omitempty"`
}

// BeforeCreate hook
func (n *UserNote) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (UserNote) TableName() string {
	return "user_notes"
}
