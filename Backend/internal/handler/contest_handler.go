package handler

import (
	"dojo/internal/dto"
	"dojo/internal/service"
	"dojo/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type ContestHandler struct {
	contestService *service.ContestService
}

func NewContestHandler(contestService *service.ContestService) *ContestHandler {
	return &ContestHandler{
		contestService: contestService,
	}
}

// ListContests handles GET /api/contests
func (h *ContestHandler) ListContests(c *fiber.Ctx) error {
	var filters dto.ContestFilterRequest

	// Parse query parameters
	filters.Platform = c.Query("platform")
	filters.Upcoming = c.QueryBool("upcoming")
	filters.Ongoing = c.QueryBool("ongoing")
	filters.Page = c.QueryInt("page", 1)
	filters.Limit = c.QueryInt("limit", 20)

	contests, total, err := h.contestService.ListContests(&filters)
	if err != nil {
		return utils.SendInternalError(c, "Failed to fetch contests", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Contests fetched successfully", fiber.Map{
		"contests": contests,
		"total":    total,
		"page":     filters.Page,
		"limit":    filters.Limit,
	})
}

// GetContest handles GET /api/contests/:id
func (h *ContestHandler) GetContest(c *fiber.Ctx) error {
	id := c.Params("id")

	contest, err := h.contestService.GetContestByID(id)
	if err != nil {
		if err == utils.ErrContestNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Contest not found", err)
		}
		return utils.SendInternalError(c, "Failed to fetch contest", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Contest fetched successfully", fiber.Map{
		"contest": contest,
	})
}

// CreateReminder handles POST /api/contests/reminders
func (h *ContestHandler) CreateReminder(c *fiber.Ctx) error {
	var req dto.CreateReminderRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	reminder, err := h.contestService.CreateReminder(userID, &req)
	if err != nil {
		if err == utils.ErrContestNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Contest not found", err)
		}
		if err.Error() == "reminder already exists for this contest" {
			return utils.SendConflict(c, "Reminder already exists for this contest")
		}
		return utils.SendInternalError(c, "Failed to create reminder", err)
	}

	return utils.SendCreated(c, "Reminder created successfully", fiber.Map{
		"reminder": reminder,
	})
}

// DeleteReminder handles DELETE /api/contests/reminders/:id
func (h *ContestHandler) DeleteReminder(c *fiber.Ctx) error {
	reminderID := c.Params("id")

	// Get user ID from context
	userID := c.Locals("userID").(string)

	err := h.contestService.DeleteReminder(userID, reminderID)
	if err != nil {
		if err == utils.ErrReminderNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Reminder not found", err)
		}
		if err == utils.ErrUnauthorized {
			return utils.SendUnauthorized(c, "You don't have permission to delete this reminder")
		}
		return utils.SendInternalError(c, "Failed to delete reminder", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Reminder deleted successfully", nil)
}
