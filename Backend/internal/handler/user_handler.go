package handler

import (
	"dojo/internal/dto"
	"dojo/internal/middleware"
	"dojo/internal/service"
	"dojo/internal/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile - GET /api/users/profile
// Retrieves the authenticated user's profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Get user ID from auth middleware
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Get user profile
	user, err := h.userService.GetProfile(userID.String())
	if err != nil {
		if err == utils.ErrUserNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "User not found", err)
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch profile", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile - PUT /api/users/profile
// Updates the authenticated user's profile information
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	// Get user ID from auth middleware
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Parse request body
	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Update profile
	user, err := h.userService.UpdateProfile(userID.String(), &req)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "User not found", err)
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update profile", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Profile updated successfully", user)
}

// UpdateUser - PATCH /api/users/account
// Updates the authenticated user's account details (username, avatar)
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Get user ID from auth middleware
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Parse request body
	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Update user
	user, err := h.userService.UpdateUser(userID.String(), &req)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "User not found", err)
		}
		if err == utils.ErrUsernameTaken {
			return utils.SendError(c, fiber.StatusConflict, "Username already taken", err)
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update account", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Account updated successfully", user)
}

// ChangePassword - POST /api/users/change-password
// Changes the authenticated user's password
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	// Get user ID from auth middleware
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Parse request body
	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Change password
	if err := h.userService.ChangePassword(userID.String(), &req); err != nil {
		if err == utils.ErrUserNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "User not found", err)
		}
		if err.Error() == "invalid old password" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Invalid old password", err)
		}
		if err.Error() == "cannot change password for OAuth-only accounts" {
			return utils.SendError(c, fiber.StatusBadRequest, "Cannot change password for OAuth-only accounts", err)
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to change password", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Password changed successfully", nil)
}

// SyncPlatformStats - POST /api/users/sync-stats
// Syncs platform statistics from external coding platforms
func (h *UserHandler) SyncPlatformStats(c *fiber.Ctx) error {
	fmt.Println("DEBUG HANDLER: SyncPlatformStats endpoint called")

	// Get user ID from auth middleware
	userID, err := middleware.GetUserID(c)
	if err != nil {
		fmt.Printf("DEBUG HANDLER: Auth error: %v\n", err)
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	fmt.Printf("DEBUG HANDLER: UserID: %s\n", userID)

	// Parse request body
	var req struct {
		Platforms []string `json:"platforms" validate:"required,min=1"`
	}
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("DEBUG HANDLER: Body parse error: %v\n", err)
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	fmt.Printf("DEBUG HANDLER: Platforms requested: %v\n", req.Platforms)

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		fmt.Printf("DEBUG HANDLER: Validation error: %v\n", err)
		return utils.SendError(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Validate platforms
	validPlatforms := map[string]bool{
		"leetcode":   true,
		"codeforces": true,
		"codechef":   true,
		"gfg":        true,
	}

	for _, platform := range req.Platforms {
		if !validPlatforms[platform] {
			return utils.SendError(c, fiber.StatusBadRequest, "Invalid platform: "+platform, nil)
		}
	}

	// Sync platform stats
	results, err := h.userService.SyncPlatformStats(userID.String(), req.Platforms)
	if err != nil {
		fmt.Printf("DEBUG HANDLER: Sync service error: %v\n", err)
		if err == utils.ErrUserNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "User not found", err)
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to sync platform stats", err)
	}

	fmt.Printf("DEBUG HANDLER: Sync results: %+v\n", results)
	return utils.SendSuccess(c, fiber.StatusOK, "Platform stats sync completed", results)
}
