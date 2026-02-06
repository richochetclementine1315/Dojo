package middleware

import (
	"strings"

	"dojo/internal/config"
	"dojo/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT token
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string

		// Get Authorization header first
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			// Extract token from header
			var err error
			token, err = utils.ExtractTokenFromHeader(authHeader)
			if err != nil {
				return utils.SendUnauthorized(c, "Invalid authorization header format")
			}
		} else {
			// Check for token in query parameter (for WebSocket connections)
			token = c.Query("token")
			if token == "" {
				return utils.SendUnauthorized(c, "Missing authorization token")
			}
		}

		// Validate token
		claims, err := utils.ValidateToken(token, cfg.JWT.Secret)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				return utils.SendUnauthorized(c, "Token has expired")
			}
			return utils.SendUnauthorized(c, "Invalid token")
		}

		// Set user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

// GetUserID gets user ID from context
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID := c.Locals("userID")
	if userID == nil {
		return uuid.Nil, utils.ErrUnauthorized
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, utils.ErrUnauthorized
	}

	return id, nil
}

// GetUserEmail gets user email from context
func GetUserEmail(c *fiber.Ctx) (string, error) {
	email := c.Locals("email")
	if email == nil {
		return "", utils.ErrUnauthorized
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", utils.ErrUnauthorized
	}

	return emailStr, nil
}
