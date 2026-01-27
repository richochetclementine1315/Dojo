package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a global error handler
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500
	code := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log error
	log.Printf("Error: %v", err)

	// Send error response
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
		"message": "An error occurred",
	})
}
