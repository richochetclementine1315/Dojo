package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserProfile struct {
	ID       uuid.UUID `gorm:"type:uuid; primaryKey; default:gen_random_uuid()" json:"id"`
	UserID   uuid.UUID `gorm:"type:uuid; not null; uniqueIndex" json:"user_id"`
	Bio      string    `gorm:"type:text" json:"bio"`
	Location string    `gorm:"type:varchar(225)" json:"location"`
	Website  string    `gorm:"type:varchar(500)" json:"website"`

	// Platform-specific usernames
	LeetcodeUsername   string `gorm:"type:varchar(100)" json:"leetcode_username"`
	CodeforcesUsername string `gorm:"type:varchar(100)" json:"codeforces_username"`
	CodechefUsername   string `gorm:"type:varchar(100)" json:"codechef_username"`
	GFGUsername        string `gorm:"type:varchar(100)" json:"gfg_username"`

	// Aggregated stats
	TotalSolved  int `gorm:"default:0" json:"total_solved"`
	EasySolved   int `gorm:"default:0" json:"easy_solved"`
	MediumSolved int `gorm:"default:0" json:"medium_solved"`
	HardSolved   int `gorm:"default:0" json:"hard_solved"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship with User
	User          User               `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
	PlatformStats []UserPlatformStat `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate Hook
func (p *UserProfile) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the UserProfile model.
func (UserProfile) TableName() string {
	return "user_profiles"
}

// UserPlatformStat holds detailed statistics for each coding platform.
type UserPlatformStat struct {
	ID            uuid.UUID  `gorm:"type:uuid; primaryKey; default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid; not null; index" json:"user_id"`
	Platform      string     `gorm:"type:varchar(50); not null" json:"platform"` // e.g., "leetcode", "codeforces"
	Rating        int        ` json:"rating"`
	MaxRating     int        ` json:"max_rating"`
	SolvedCount   int        `gorm:"default:0" json:"solved_count"`
	ContestRating int        ` json:"contest_rating"`
	GlobalRank    int        ` json:"global_rank"`
	CountryRank   int        ` json:"country_rank"`
	RawData       string     `gorm:"type:jsonb" json:"raw_data"` // Store complete API response or additional data
	LastSyncedAt  *time.Time ` json:"last_synced_at"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship with User
	User User `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate Hook
func (s *UserPlatformStat) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the UserPlatformStat model.
func (UserPlatformStat) TableName() string {
	return "user_platform_stats"
}
