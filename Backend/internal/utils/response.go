package utils

import (
	"github.com/gofiber/fiber/v2"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// SendSuccess sends a success response
func SendSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendCreated sends a created (201) response
func SendCreated(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends an error response
func SendError(c *fiber.Ctx, statusCode int, err error, message string) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	return c.Status(statusCode).JSON(ErrorResponse{
		Success: false,
		Error:   errorMsg,
		Message: message,
	})
}

// SendBadRequest sends a 400 Bad Request response
func SendBadRequest(c *fiber.Ctx, err error, message string) error {
	return SendError(c, fiber.StatusBadRequest, err, message)
}

// SendUnauthorized sends a 401 Unauthorized response
func SendUnauthorized(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusUnauthorized, nil, message)
}

// SendForbidden sends a 403 Forbidden response
func SendForbidden(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusForbidden, nil, message)
}

// SendNotFound sends a 404 Not Found response
func SendNotFound(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusNotFound, nil, message)
}

// SendConflict sends a 409 Conflict response
func SendConflict(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusConflict, nil, message)
}

// SendInternalError sends a 500 Internal Server Error response
func SendInternalError(c *fiber.Ctx, err error, message string) error {
	return SendError(c, fiber.StatusInternalServerError, err, message)
}

// SendPaginated sends a paginated response
func SendPaginated(c *fiber.Ctx, data interface{}, page, limit int, total int64) error {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
