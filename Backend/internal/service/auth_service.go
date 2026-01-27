package service

import (
	"context"
	"fmt"
	"time"

	"dojo/internal/config"
	"dojo/internal/dto"
	"dojo/internal/models"
	"dojo/internal/repository"
	"dojo/internal/utils"
	"dojo/pkg/oauth"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
	authRepo *repository.AuthRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, authRepo *repository.AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		authRepo: authRepo,
		cfg:      cfg,
	}
}

// Register registers a new user with email/password
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.TokenResponse, error) {
	// Check if email exists
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.ErrEmailTaken
	}

	// Check if username exists
	exists, err = s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.ErrUsernameTaken
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		IsVerified:   false,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Create profile
	profile := &models.UserProfile{
		UserID:             user.ID,
		LeetcodeUsername:   req.LeetcodeUsername,
		CodeforcesUsername: req.CodeforcesUsername,
		CodechefUsername:   req.CodechefUsername,
		GFGUsername:        req.GFGUsername,
	}

	if err := s.userRepo.CreateProfile(profile); err != nil {
		return nil, err
	}

	// Create auth account
	authAccount := &models.AuthAccount{
		UserID:         user.ID,
		Provider:       "email",
		ProviderUserID: user.Email,
	}

	if err := s.authRepo.CreateAuthAccount(authAccount); err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateTokens(user)
}

// Login authenticates user with email/password
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.TokenResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user has password (not OAuth-only)
	if user.PasswordHash == "" {
		return nil, utils.ErrInvalidCredentials
	}

	// Verify password
	isValid := utils.ComparePassword(user.PasswordHash, req.Password)
	if !isValid {
		return nil, utils.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate tokens
	return s.generateTokens(user)
}

// GoogleOAuth handles Google OAuth login
func (s *AuthService) GoogleOAuth(code string) (*dto.TokenResponse, error) {
	ctx := context.Background()

	// Setup OAuth config
	oauthConfig := oauth.GoogleOAuthConfig(
		s.cfg.OAuth.Google.ClientID,
		s.cfg.OAuth.Google.ClientSecret,
		s.cfg.OAuth.Google.RedirectURL,
	)

	// Exchange code for token
	token, err := oauth.ExchangeGoogleCode(ctx, oauthConfig, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info
	googleUser, err := oauth.GetGoogleUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Find or create user
	return s.findOrCreateOAuthUser("google", googleUser.ID, googleUser.Email, googleUser.Name, googleUser.Picture, token)
}

// GitHubOAuth handles GitHub OAuth login
func (s *AuthService) GitHubOAuth(code string) (*dto.TokenResponse, error) {
	ctx := context.Background()

	// Setup OAuth config
	oauthConfig := oauth.GitHubOAuthConfig(
		s.cfg.OAuth.GitHub.ClientID,
		s.cfg.OAuth.GitHub.ClientSecret,
		s.cfg.OAuth.GitHub.RedirectURL,
	)

	// Exchange code for token
	token, err := oauth.ExchangeGitHubCode(ctx, oauthConfig, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info
	githubUser, err := oauth.GetGitHubUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	username := githubUser.Login
	if githubUser.Name != "" {
		username = githubUser.Name
	}

	// Find or create user
	return s.findOrCreateOAuthUser("github", fmt.Sprintf("%d", githubUser.ID), githubUser.Email, username, githubUser.AvatarURL, token)
}

// RefreshAccessToken refreshes access token using refresh token
func (s *AuthService) RefreshAccessToken(refreshTokenStr string) (*dto.TokenResponse, error) {
	// Find refresh token
	refreshToken, err := s.authRepo.FindRefreshToken(refreshTokenStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrInvalidToken
		}
		return nil, err
	}

	// Check if expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		s.authRepo.DeleteRefreshToken(refreshTokenStr)
		return nil, utils.ErrTokenExpired
	}

	// Generate new tokens
	return s.generateTokens(&refreshToken.User)
}

// Logout logs out user by deleting refresh token
func (s *AuthService) Logout(refreshToken string) error {
	return s.authRepo.DeleteRefreshToken(refreshToken)
}

// Helper: findOrCreateOAuthUser finds or creates user from OAuth
func (s *AuthService) findOrCreateOAuthUser(provider, providerUserID, email, username, avatarURL string, token interface{}) (*dto.TokenResponse, error) {
	// Check if auth account exists
	authAccount, err := s.authRepo.FindAuthAccount(provider, providerUserID)
	if err == nil {
		// User exists, return tokens
		return s.generateTokens(&authAccount.User)
	}

	// User doesn't exist, create new
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil {
		// Email exists, link OAuth account
		newAuthAccount := &models.AuthAccount{
			UserID:         existingUser.ID,
			Provider:       provider,
			ProviderUserID: providerUserID,
		}
		if err := s.authRepo.CreateAuthAccount(newAuthAccount); err != nil {
			return nil, err
		}
		return s.generateTokens(existingUser)
	}

	// Create new user
	user := &models.User{
		Email:      email,
		Username:   s.generateUniqueUsername(username),
		AvatarURL:  avatarURL,
		IsVerified: true, // OAuth users are verified
		IsActive:   true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Create profile
	profile := &models.UserProfile{
		UserID: user.ID,
	}
	if err := s.userRepo.CreateProfile(profile); err != nil {
		return nil, err
	}

	// Create auth account
	newAuthAccount := &models.AuthAccount{
		UserID:         user.ID,
		Provider:       provider,
		ProviderUserID: providerUserID,
	}
	if err := s.authRepo.CreateAuthAccount(newAuthAccount); err != nil {
		return nil, err
	}

	return s.generateTokens(user)
}

// Helper: generateTokens generates access and refresh tokens
func (s *AuthService) generateTokens(user *models.User) (*dto.TokenResponse, error) {
	// Generate access token
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
		s.cfg.JWT.Secret,
		s.cfg.JWT.AccessExpiry,
	)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenStr, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Save refresh token to database
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshExpiry),
	}

	if err := s.authRepo.CreateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.cfg.JWT.AccessExpiry.Seconds()),
	}, nil
}

// Helper: generateUniqueUsername generates a unique username
func (s *AuthService) generateUniqueUsername(baseUsername string) string {
	username := baseUsername
	counter := 1

	for {
		exists, _ := s.userRepo.ExistsByUsername(username)
		if !exists {
			return username
		}
		username = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++
	}
}
