package handler

import (
	"dojo/internal/dto"
	"dojo/internal/middleware"
	"dojo/internal/service"
	"dojo/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProblemHandler struct {
	problemService *service.ProblemService
}

func NewProblemHandler(problemService *service.ProblemService) *ProblemHandler {
	return &ProblemHandler{
		problemService: problemService,
	}
}

// CreateProblem - POST /api/problems
// Creates a new problem (admin only)
func (h *ProblemHandler) CreateProblem(c *fiber.Ctx) error {
	var req dto.CreateProblemRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request body", err)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, "Validation failed", err)
	}

	problem, err := h.problemService.CreateProblem(&req)
	if err != nil {
		if err.Error() == "problem with this URL already exists" {
			return utils.SendConflict(c, err.Error())
		}
		return utils.SendInternalError(c, "Failed to create problem", err)
	}

	return utils.SendCreated(c, "Problem created successfully", problem)
}

// GetProblem - GET /api/problems/:id
// Retrieves a problem by ID
func (h *ProblemHandler) GetProblem(c *fiber.Ctx) error {
	id := c.Params("id")

	problem, err := h.problemService.GetProblemByID(id)
	if err != nil {
		if err == utils.ErrProblemNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Problem not found", err)
		}
		return utils.SendInternalError(c, "Failed to fetch problem", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Problem retrieved successfully", problem)
}

// ListProblems - GET /api/problems
// Lists all problems with filters and pagination
func (h *ProblemHandler) ListProblems(c *fiber.Ctx) error {
	var filters dto.ProblemFilterRequest
	if err := c.QueryParser(&filters); err != nil {
		return utils.SendBadRequest(c, "Invalid query parameters", err)
	}

	// Get user ID if authenticated (optional)
	userID, _ := middleware.GetUserID(c)
	var userIDStr string
	if userID != uuid.Nil {
		userIDStr = userID.String()
	}

	problems, total, err := h.problemService.ListProblems(&filters, userIDStr)
	if err != nil {
		return utils.SendInternalError(c, "Failed to fetch problems", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Problems retrieved successfully", fiber.Map{
		"problems": problems,
		"total":    total,
		"page":     filters.Page,
		"limit":    filters.Limit,
	})
}

// UpdateProblem - PUT /api/problems/:id
// Updates a problem (admin only)
func (h *ProblemHandler) UpdateProblem(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateProblemRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request body", err)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, "Validation failed", err)
	}

	problem, err := h.problemService.UpdateProblem(id, &req)
	if err != nil {
		if err == utils.ErrProblemNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Problem not found", err)
		}
		return utils.SendInternalError(c, "Failed to update problem", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Problem updated successfully", problem)
}

// DeleteProblem - DELETE /api/problems/:id
// Deletes a problem (admin only)
func (h *ProblemHandler) DeleteProblem(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.problemService.DeleteProblem(id); err != nil {
		if err == utils.ErrProblemNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Problem not found", err)
		}
		return utils.SendInternalError(c, "Failed to delete problem", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Problem deleted successfully", nil)
}

// SyncProblems - POST /api/problems/sync
// Syncs problems from external platforms
func (h *ProblemHandler) SyncProblems(c *fiber.Ctx) error {
	var req struct {
		Platform string `json:"platform" validate:"required,oneof=leetcode codeforces"`
		Limit    int    `json:"limit"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request body", err)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, "Validation failed", err)
	}

	if req.Limit <= 0 {
		req.Limit = 100 // Default limit
	}

	count, err := h.problemService.SyncProblems(req.Platform, req.Limit)
	if err != nil {
		return utils.SendInternalError(c, "Failed to sync problems", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Problems synced successfully", fiber.Map{
		"imported": count,
		"platform": req.Platform,
	})
}

// MarkProblemSolved - POST /api/problems/:id/solve
// Marks a problem as solved or unsolved for the authenticated user
func (h *ProblemHandler) MarkProblemSolved(c *fiber.Ctx) error {
	problemID := c.Params("id")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	var req struct {
		IsSolved bool `json:"is_solved"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request body", err)
	}

	if err := h.problemService.MarkProblemSolved(userID.String(), problemID, req.IsSolved); err != nil {
		return utils.SendInternalError(c, "Failed to update problem status", err)
	}

	message := "Problem marked as solved"
	if !req.IsSolved {
		message = "Problem marked as unsolved"
	}

	return utils.SendSuccess(c, fiber.StatusOK, message, nil)
}

// GetUserSolvedCount - GET /api/problems/solved/count
// Returns the count of problems solved by the authenticated user
func (h *ProblemHandler) GetUserSolvedCount(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	count, err := h.problemService.GetUserSolvedCount(userID.String())
	if err != nil {
		return utils.SendInternalError(c, "Failed to get solved count", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Solved count retrieved successfully", fiber.Map{
		"count": count,
	})
}
