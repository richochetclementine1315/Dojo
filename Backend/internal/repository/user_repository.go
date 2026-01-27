package repository

import (
	"dojo/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID retrieves a user by ID

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = r.db.Preload("Profile").First(&user, "id=?", userID).Error
	return &user, err
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Profile").First(&user, "email=?", email).Error
	return &user, err
}

// FindByUsername retrieves a user by username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Profile").First(&user, "username=?", username).Error
	return &user, err
}

// Upadate updates an existing user in the database
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// ExistsByEmail checks if a user exists with the given email
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email=?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByUsername checks if a user exists with the given username
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username=?", username).Count(&count).Error
	return count > 0, err
}

// CreateProfile creates a new user profile in the database
func (r *UserRepository) CreateProfile(profile *models.UserProfile) error {
	return r.db.Create(profile).Error
}

// GetProfile gets a user's profile from the database
func (r *UserRepository) GetProfile(userID uuid.UUID) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := r.db.Preload("PlatformStats").First(&profile, "user_id=?", userID).Error
	return &profile, err
}

// Add these methods to user_repository.go after GetProfile:

// LoadProfile loads the user's profile relationship
func (r *UserRepository) LoadProfile(user *models.User) error {
	return r.db.Preload("Profile").First(user, "id = ?", user.ID).Error
}

// LoadPlatformStats loads the user's platform stats relationship
func (r *UserRepository) LoadPlatformStats(user *models.User) error {
	return r.db.Preload("PlatformStats").First(user, "id = ?", user.ID).Error
}

// UpdateProfile updates an existing user profile
func (r *UserRepository) UpdateProfile(profile *models.UserProfile) error {
	return r.db.Save(profile).Error
}

// UpsertPlatformStat creates or updates a platform stat
func (r *UserRepository) UpsertPlatformStat(stat *models.UserPlatformStat) error {
	// Check if stat exists for this user and platform
	var existing models.UserPlatformStat
	err := r.db.Where("user_id = ? AND platform = ?", stat.UserID, stat.Platform).First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new stat
			return r.db.Create(stat).Error
		}
		return err
	}

	// Update existing stat
	stat.ID = existing.ID
	return r.db.Save(stat).Error
}
