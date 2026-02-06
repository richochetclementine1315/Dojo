package handler

import (
	"dojo/internal/config"
	"dojo/internal/dto"
	"dojo/internal/service"
	"dojo/internal/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
	config      *config.Config
}

func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      cfg,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid Request Body", err)
	}
	// validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, "Validation Error", err)
	}
	// Register user
	tokenResponse, err := h.authService.Register(&req)
	if err != nil {
		if err == utils.ErrEmailTaken || err == utils.ErrUsernameTaken {
			return utils.SendConflict(c, err.Error())
		}
		return utils.SendInternalError(c, "Failed to register user", err)
	}
	return utils.SendCreated(c, "User registered successfully", tokenResponse)
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid Request Body", err)
	}
	// validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, "Validation Failed", err)
	}
	// Login User
	tokenResponse, err := h.authService.Login(&req)
	if err != nil {
		if err == utils.ErrInvalidCredentials {
			return utils.SendUnauthorized(c, "Invalid Email or password")
		}
		return utils.SendInternalError(c, "Failed to login user", err)
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Login Successful", tokenResponse)
}

// GoogleLogin handles Google OAuth login
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	// Build Google OAuth URL with parameters
	url := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile&access_type=offline",
		h.config.OAuth.Google.ClientID,
		h.config.OAuth.Google.RedirectURL,
	)
	// Redirect browser to Google
	return c.Redirect(url, fiber.StatusFound)
}

// GoogleCallback handles Google OAuth callback
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return utils.SendBadRequest(c, "Authorization Code required", nil)
	}
	// Handle OAuth
	tokenResponse, err := h.authService.GoogleOAuth(code)
	if err != nil {
		return utils.SendInternalError(c, "Failed to authenticate with Google", err)
	}

	// Redirect to frontend callback with tokens
	redirectURL := fmt.Sprintf(
		"http://localhost:5173/auth/google/callback?access_token=%s&refresh_token=%s",
		tokenResponse.AccessToken,
		tokenResponse.RefreshToken,
	)
	return c.Redirect(redirectURL, fiber.StatusFound)
}

// GithubLogin handles GitHub OAuth login
func (h *AuthHandler) GitHubLogin(c *fiber.Ctx) error {
	// Build GitHub OAuth URL with parameters
	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		h.config.OAuth.GitHub.ClientID,
		h.config.OAuth.GitHub.RedirectURL,
	)
	// Redirect browser to GitHub
	return c.Redirect(url, fiber.StatusFound)
}

// GitHubCallback handles GitHub OAuth callback
func (h *AuthHandler) GitHubCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return utils.SendBadRequest(c, "Authorization Code is required", nil)
	}
	// Handle OAuth
	tokenResponse, err := h.authService.GitHubOAuth(code)
	if err != nil {
		return utils.SendInternalError(c, "Failed to authenticate with GitHub", err)
	}

	// Redirect to frontend callback with tokens
	redirectURL := fmt.Sprintf(
		"http://localhost:5173/auth/github/callback?access_token=%s&refresh_token=%s",
		tokenResponse.AccessToken,
		tokenResponse.RefreshToken,
	)
	return c.Redirect(redirectURL, fiber.StatusFound)
}

// RefreshToken handles token refreshing
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request body", err)
	}
	// validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, "Validation failed", err)
	}
	// Refresh Token
	tokenResponse, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		if err == utils.ErrInvalidToken || err == utils.ErrTokenExpired {
			return utils.SendUnauthorized(c, err.Error())
		}
		return utils.SendInternalError(c, "Failed to refresh token", err)
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Token refreshed successfully", tokenResponse)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest

	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request body", err)
	}
	// Logout User
	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		return utils.SendInternalError(c, "Failed to logout user", err)
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Logout successful", nil)
}
