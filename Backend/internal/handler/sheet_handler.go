package handler

import (
	"dojo/internal/dto"
	"dojo/internal/service"
	"dojo/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SheetHandler struct {
	sheetService *service.SheetService
}

func NewSheetHandler(sheetService *service.SheetService) *SheetHandler {
	return &SheetHandler{
		sheetService: sheetService,
	}
}

// CreateSheet handles POST /api/sheets
func (h *SheetHandler) CreateSheet(c *fiber.Ctx) error {
	var req dto.CreateSheetRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	userID := c.Locals("userID").(uuid.UUID).String()

	sheet, err := h.sheetService.CreateSheet(userID, &req)
	if err != nil {
		return utils.SendInternalError(c, "Failed to create sheet", err)
	}

	return utils.SendCreated(c, "Sheet created successfully", fiber.Map{
		"sheet": sheet,
	})
}

// GetSheet handles GET /api/sheets/:id
func (h *SheetHandler) GetSheet(c *fiber.Ctx) error {
	sheetID := c.Params("id")
	userID := c.Locals("userID").(uuid.UUID).String()

	sheet, err := h.sheetService.GetSheetByID(userID, sheetID)
	if err != nil {
		if err == utils.ErrSheetNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Sheet not found", err)
		}
		if err == utils.ErrSheetAccessDenied {
			return utils.SendError(c, fiber.StatusForbidden, "Access denied to this sheet", err)
		}
		return utils.SendInternalError(c, "Failed to fetch sheet", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Sheet fetched successfully", fiber.Map{
		"sheet": sheet,
	})
}

// GetUserSheets handles GET /api/sheets
func (h *SheetHandler) GetUserSheets(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID).String()

	sheets, err := h.sheetService.GetUserSheets(userID)
	if err != nil {
		return utils.SendInternalError(c, "Failed to fetch sheets", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Sheets fetched successfully", fiber.Map{
		"sheets": sheets,
	})
}

// GetPublicSheets handles GET /api/sheets/public
func (h *SheetHandler) GetPublicSheets(c *fiber.Ctx) error {
	sheets, err := h.sheetService.GetPublicSheets()
	if err != nil {
		return utils.SendInternalError(c, "Failed to fetch public sheets", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Public sheets fetched successfully", fiber.Map{
		"sheets": sheets,
	})
}

// UpdateSheet handles PUT /api/sheets/:id
func (h *SheetHandler) UpdateSheet(c *fiber.Ctx) error {
	sheetID := c.Params("id")
	userID := c.Locals("userID").(uuid.UUID).String()

	var req dto.UpdateSheetRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	sheet, err := h.sheetService.UpdateSheet(userID, sheetID, &req)
	if err != nil {
		if err == utils.ErrSheetNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Sheet not found", err)
		}
		if err == utils.ErrSheetAccessDenied {
			return utils.SendError(c, fiber.StatusForbidden, "Access denied to this sheet", err)
		}
		return utils.SendInternalError(c, "Failed to update sheet", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Sheet updated successfully", fiber.Map{
		"sheet": sheet,
	})
}

// DeleteSheet handles DELETE /api/sheets/:id
func (h *SheetHandler) DeleteSheet(c *fiber.Ctx) error {
	sheetID := c.Params("id")
	userID := c.Locals("userID").(uuid.UUID).String()

	err := h.sheetService.DeleteSheet(userID, sheetID)
	if err != nil {
		if err == utils.ErrSheetNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Sheet not found", err)
		}
		if err == utils.ErrSheetAccessDenied {
			return utils.SendError(c, fiber.StatusForbidden, "Access denied to this sheet", err)
		}
		return utils.SendInternalError(c, "Failed to delete sheet", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Sheet deleted successfully", nil)
}

// AddProblemToSheet handles POST /api/sheets/:id/problems
func (h *SheetHandler) AddProblemToSheet(c *fiber.Ctx) error {
	sheetID := c.Params("id")
	userID := c.Locals("userID").(uuid.UUID).String()

	var req dto.AddProblemToSheetRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	problem, err := h.sheetService.AddProblemToSheet(userID, sheetID, &req)
	if err != nil {
		if err == utils.ErrSheetNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Sheet not found", err)
		}
		if err == utils.ErrProblemNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Problem not found", err)
		}
		if err == utils.ErrSheetAccessDenied {
			return utils.SendError(c, fiber.StatusForbidden, "Access denied to this sheet", err)
		}
		if err == utils.ErrProblemAlreadyInSheet {
			return utils.SendError(c, fiber.StatusConflict, "Problem already exists in this sheet", err)
		}
		return utils.SendInternalError(c, "Failed to add problem to sheet", err)
	}

	return utils.SendCreated(c, "Problem added to sheet successfully", fiber.Map{
		"problem": problem,
	})
}

// RemoveProblemFromSheet handles DELETE /api/sheets/:id/problems/:problemId
func (h *SheetHandler) RemoveProblemFromSheet(c *fiber.Ctx) error {
	sheetID := c.Params("id")
	problemID := c.Params("problemId")
	userID := c.Locals("userID").(uuid.UUID).String()

	err := h.sheetService.RemoveProblemFromSheet(userID, sheetID, problemID)
	if err != nil {
		if err == utils.ErrSheetNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Sheet not found", err)
		}
		if err == utils.ErrSheetAccessDenied {
			return utils.SendError(c, fiber.StatusForbidden, "Access denied to this sheet", err)
		}
		return utils.SendInternalError(c, "Failed to remove problem from sheet", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Problem removed from sheet successfully", nil)
}

// UpdateSheetProblem handles PATCH /api/sheets/:id/problems/:problemId
func (h *SheetHandler) UpdateSheetProblem(c *fiber.Ctx) error {
	sheetID := c.Params("id")
	problemID := c.Params("problemId")
	userID := c.Locals("userID").(uuid.UUID).String()

	var req dto.UpdateSheetProblemRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	problem, err := h.sheetService.UpdateSheetProblem(userID, sheetID, problemID, &req)
	if err != nil {
		if err == utils.ErrSheetNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Sheet not found", err)
		}
		if err == utils.ErrProblemNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Problem not found in sheet", err)
		}
		if err == utils.ErrSheetAccessDenied {
			return utils.SendError(c, fiber.StatusForbidden, "Access denied to this sheet", err)
		}
		return utils.SendInternalError(c, "Failed to update problem", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Problem updated successfully", fiber.Map{
		"problem": problem,
	})
}
