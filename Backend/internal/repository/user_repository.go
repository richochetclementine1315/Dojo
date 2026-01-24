package repository

import (
	"dojo/internal/models"

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
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Profile").First(&user, "id=?", id).Error
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
