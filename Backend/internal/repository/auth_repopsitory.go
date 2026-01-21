package repository

import (
	"time"

	"dojo/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateAuthAccount creates a new auth account
func (r *AuthRepository) CreateAuthAccount(account *models.AuthAccount) error {
	return r.db.Create(account).Error
}

// FindAuthAccount finds an auth account by provider and provider user ID
func (r *AuthRepository) FindAuthAccount(provider, providerUserID string) (*models.AuthAccount, error) {
	var account models.AuthAccount
	err := r.db.Preload("User").First(&account, "provider = ? AND provider_user_id = ?", provider, providerUserID).Error
	return &account, err
}

// UpdateAuthAccount updates an auth account
func (r *AuthRepository) UpdateAuthAccount(account *models.AuthAccount) error {
	return r.db.Save(account).Error
}

// CreateRefreshToken creates a new refresh token
func (r *AuthRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

// FindRefreshToken finds a refresh token by token string
func (r *AuthRepository) FindRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Preload("User").First(&refreshToken, "token = ?", token).Error
	return &refreshToken, err
}

// DeleteRefreshToken deletes a refresh token
func (r *AuthRepository) DeleteRefreshToken(token string) error {
	return r.db.Delete(&models.RefreshToken{}, "token = ?", token).Error
}

// DeleteExpiredTokens deletes all expired refresh tokens
func (r *AuthRepository) DeleteExpiredTokens() error {
	return r.db.Delete(&models.RefreshToken{}, "expires_at < ?", time.Now()).Error
}

// DeleteUserTokens deletes all refresh tokens for a user (logout)
func (r *AuthRepository) DeleteUserTokens(userID uuid.UUID) error {
	return r.db.Delete(&models.RefreshToken{}, "user_id = ?", userID).Error
}
