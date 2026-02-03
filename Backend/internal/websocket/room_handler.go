package websocket

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RoomHandler handles WebSocket connections for rooms
type RoomHandler struct {
	Hub *Hub
}

// NewRoomHandler creates a new RoomHandler
func NewRoomHandler(hub *Hub) *RoomHandler {
	return &RoomHandler{
		Hub: hub,
	}
}

// UpgradeConnection upgrades HTTP to WebSocket for room connection
func (h *RoomHandler) UpgradeConnection(c *fiber.Ctx) error {
	// Get room ID from URL params
	roomIDStr := c.Params("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	// Get user info from auth middleware (should be set by middleware)
	userUUID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get email for username (we don't have username in JWT claims)
	email, ok := c.Locals("email").(string)
	if !ok {
		email = "Anonymous"
	}

	// Check if request is WebSocket upgrade
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		c.Locals("roomID", roomID)
		c.Locals("userUUID", userUUID)
		c.Locals("email", email)
		return c.Next()
	}

	return c.Status(fiber.StatusUpgradeRequired).JSON(fiber.Map{
		"error": "WebSocket upgrade required",
	})
}

// HandleConnection handles the WebSocket connection
func (h *RoomHandler) HandleConnection(c *websocket.Conn) {
	// Get metadata from locals
	roomID := c.Locals("roomID").(uuid.UUID)
	userUUID := c.Locals("userUUID").(uuid.UUID)
	email := c.Locals("email").(string)

	// Create client
	client := NewClient(userUUID, email, roomID, c, h.Hub)

	// Register client
	h.Hub.Register <- client

	log.Printf("WebSocket connected: User %s in room %s", email, roomID)

	// Start read and write pumps
	go client.WritePump()
	client.ReadPump() // Blocking call
}
