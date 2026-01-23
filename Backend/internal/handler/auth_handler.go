package handler

import (
	"dojo/internal/dto"
	"dojo/internal/service"
	"dojo/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, err, "Invalid Request Body")
	}
	// validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, err, "Validation Error")
	}
	// Register user
	tokenResponse, err := h.authService.Register(&req)
	if err != nil {
		if err == utils.ErrEmailTaken || err == utils.ErrUsernameTaken {
			return utils.SendConflict(c, err.Error())
		}
		return utils.SendInternalError(c, err, "Failed to register user")
	}
	return utils.SendCreated(c, "User registered successfully", tokenResponse)
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, err, "Invalid Request Body")
	}
	// validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, err, "Validation Failed")
	}
	// Login User
	tokenResponse, err := h.authService.Login(&req)
	if err != nil {
		if err == utils.ErrInvalidCredentials {
			return utils.SendUnauthorized(c, "Invalid Email or password")
		}
		return utils.SendInternalError(c, err, "Failed to login user")
	}
	return utils.SendSuccess(c, "Login Successful", tokenResponse)
}

// GoogleLogin handles Google OAuth login
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	// Return OAuth URL for the frontend to redirect
	url := "https://accounts.google.com/o/oauth2/v2/auth"
	return utils.SendSuccess(c, "Google OAuth URL", fiber.Map{
		"url":     url,
		"message": "Redirect User to this URL with appropriate parameters",
	})
}

// GoogleCallback handles Google OAuth callback
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return utils.SendBadRequest(c, nil, "Authorization Code required")
	}
	// Handle OAuth
	tokenResponse, err := h.authService.GoogleOAuth(code)
	if err != nil {
		return utils.SendInternalError(c, err, "Failed to authenticate with Google")
	}
	return utils.SendSuccess(c, "Google authentication successful", tokenResponse)
}

// GithubLogin handles GitHub OAuth login
func (h *AuthHandler) GitHubLogin(c *fiber.Ctx) error {
	// Return OAuth URL for the frontend to redirect
	url := "https://github.com/login/oauth/authorize"
	return utils.SendSuccess(c, "GitHub OAuth URL", fiber.Map{
		"url":     url,
		"message": "Redirect User to this URL with appropriate parameters",
	})
}

// GitHubCallback handles GitHub OAuth callback
func (h *AuthHandler) GitHubCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return utils.SendBadRequest(c, nil, "Authorization Code is required")
	}
	// Handle OAuth
	tokenResponse, err := h.authService.GitHubOAuth(code)
	if err != nil {
		return utils.SendInternalError(c, err, "Failed to authenticate with GitHub")
	}
	return utils.SendSuccess(c, "GitHub authentication successful", tokenResponse)
}

// RefreshToken handles token refreshing
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, err, "Invalid request body")
	}
	// validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, err, "Validation failed")
	}
	// Refresh Token
	tokenResponse, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		if err == utils.ErrInvalidToken || err == utils.ErrTokenExpired {
			return utils.SendUnauthorized(c, err.Error())
		}
		return utils.SendInternalError(c, err, "Failed to refresh token")
	}
	return utils.SendSuccess(c, "Token refreshed successfully", tokenResponse)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest

	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, err, "Invalid request body")
	}
	// Logout User
	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		return utils.SendInternalError(c, err, "Failed to logout user")
	}
	return utils.SendSuccess(c, "Logout successful", nil)
}
