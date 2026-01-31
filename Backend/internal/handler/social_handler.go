package handler

import (
	"dojo/internal/dto"
	"dojo/internal/service"
	"dojo/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type SocialHandler struct {
	socialService *service.SocialService
}

func NewSocialHandler(socialService *service.SocialService) *SocialHandler {
	return &SocialHandler{
		socialService: socialService,
	}
}

// SendFriendRequest handles POST /api/social/friends/requests
func (h *SocialHandler) SendFriendRequest(c *fiber.Ctx) error {
	var req dto.SendFriendRequestRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	userID := c.Locals("userID").(string)

	friendRequest, err := h.socialService.SendFriendRequest(userID, &req)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "User not found", err)
		}
		if err == utils.ErrAlreadyFriends {
			return utils.SendError(c, fiber.StatusConflict, "Already friends with this user", err)
		}
		if err == utils.ErrCannotSendToSelf {
			return utils.SendBadRequest(c, "Cannot send friend request to yourself", err)
		}
		if err == utils.ErrFriendRequestAlreadyExists {
			return utils.SendError(c, fiber.StatusConflict, "Friend request already exists", err)
		}
		if err == utils.ErrUserBlocked {
			return utils.SendError(c, fiber.StatusForbidden, "Cannot send request to this user", err)
		}
		return utils.SendInternalError(c, "Failed to send friend request", err)
	}

	return utils.SendCreated(c, "Friend request sent successfully", fiber.Map{
		"request": friendRequest,
	})
}

// GetReceivedRequests handles GET /api/social/friends/requests/received
func (h *SocialHandler) GetReceivedRequests(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	requests, err := h.socialService.GetReceivedRequests(userID)
	if err != nil {
		return utils.SendInternalError(c, "Failed to fetch received requests", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Received requests fetched successfully", fiber.Map{
		"requests": requests,
	})
}

// GetSentRequests handles GET /api/social/friends/requests/sent
func (h *SocialHandler) GetSentRequests(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	requests, err := h.socialService.GetSentRequests(userID)
	if err != nil {
		return utils.SendInternalError(c, "Failed to fetch sent requests", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Sent requests fetched successfully", fiber.Map{
		"requests": requests,
	})
}

// RespondToFriendRequest handles PATCH /api/social/friends/requests/:id
func (h *SocialHandler) RespondToFriendRequest(c *fiber.Ctx) error {
	requestID := c.Params("id")
	userID := c.Locals("userID").(string)

	var req dto.UpdateFriendRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	err := h.socialService.RespondToFriendRequest(userID, requestID, &req)
	if err != nil {
		if err == utils.ErrFriendRequestNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Friend request not found", err)
		}
		if err == utils.ErrUnauthorized {
			return utils.SendError(c, fiber.StatusForbidden, "Not authorized to respond to this request", err)
		}
		if err.Error() == "friend request already processed" {
			return utils.SendError(c, fiber.StatusConflict, "Friend request already processed", err)
		}
		return utils.SendInternalError(c, "Failed to respond to friend request", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Friend request "+req.Action+"ed successfully", nil)
}

// CancelFriendRequest handles DELETE /api/social/friends/requests/:id
func (h *SocialHandler) CancelFriendRequest(c *fiber.Ctx) error {
	requestID := c.Params("id")
	userID := c.Locals("userID").(string)

	err := h.socialService.CancelFriendRequest(userID, requestID)
	if err != nil {
		if err == utils.ErrFriendRequestNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Friend request not found", err)
		}
		if err == utils.ErrUnauthorized {
			return utils.SendError(c, fiber.StatusForbidden, "Not authorized to cancel this request", err)
		}
		return utils.SendInternalError(c, "Failed to cancel friend request", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Friend request cancelled successfully", nil)
}

// GetFriends handles GET /api/social/friends
func (h *SocialHandler) GetFriends(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	friends, err := h.socialService.GetFriends(userID)
	if err != nil {
		return utils.SendInternalError(c, "Failed to fetch friends", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Friends fetched successfully", fiber.Map{
		"friends": friends,
	})
}

// RemoveFriend handles DELETE /api/social/friends/:id
func (h *SocialHandler) RemoveFriend(c *fiber.Ctx) error {
	friendID := c.Params("id")
	userID := c.Locals("userID").(string)

	err := h.socialService.RemoveFriend(userID, friendID)
	if err != nil {
		if err.Error() == "not friends with this user" {
			return utils.SendError(c, fiber.StatusBadRequest, "Not friends with this user", err)
		}
		return utils.SendInternalError(c, "Failed to remove friend", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Friend removed successfully", nil)
}

// BlockUser handles POST /api/social/blocks
func (h *SocialHandler) BlockUser(c *fiber.Ctx) error {
	var req dto.BlockUserRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	userID := c.Locals("userID").(string)

	err := h.socialService.BlockUser(userID, &req)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "User not found", err)
		}
		if err.Error() == "cannot block yourself" {
			return utils.SendBadRequest(c, "Cannot block yourself", err)
		}
		if err.Error() == "user already blocked" {
			return utils.SendError(c, fiber.StatusConflict, "User already blocked", err)
		}
		return utils.SendInternalError(c, "Failed to block user", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User blocked successfully", nil)
}

// UnblockUser handles DELETE /api/social/blocks/:id
func (h *SocialHandler) UnblockUser(c *fiber.Ctx) error {
	blockedID := c.Params("id")
	userID := c.Locals("userID").(string)

	err := h.socialService.UnblockUser(userID, blockedID)
	if err != nil {
		if err.Error() == "user not blocked" {
			return utils.SendError(c, fiber.StatusBadRequest, "User not blocked", err)
		}
		return utils.SendInternalError(c, "Failed to unblock user", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User unblocked successfully", nil)
}

// GetBlockedUsers handles GET /api/social/blocks
func (h *SocialHandler) GetBlockedUsers(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	users, err := h.socialService.GetBlockedUsers(userID)
	if err != nil {
		return utils.SendInternalError(c, "Failed to fetch blocked users", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Blocked users fetched successfully", fiber.Map{
		"users": users,
	})
}

// SearchUsers handles GET /api/social/users/search
func (h *SocialHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	limit := c.QueryInt("limit", 20)

	if query == "" {
		return utils.SendBadRequest(c, "Search query is required", nil)
	}

	users, err := h.socialService.SearchUsers(query, limit)
	if err != nil {
		return utils.SendInternalError(c, "Failed to search users", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Users fetched successfully", fiber.Map{
		"users": users,
	})
}
