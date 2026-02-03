package handler

import (
	"dojo/internal/dto"
	"dojo/internal/service"
	"dojo/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RoomHandler struct {
	roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

// CreateRoom handles POST /api/rooms
func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var req dto.CreateRoomRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, "Validation failed", err)
	}

	userID := c.Locals("userID").(uuid.UUID).String()

	room, err := h.roomService.CreateRoom(userID, &req)
	if err != nil {
		return utils.SendInternalError(c, "Failed to create room", err)
	}

	return utils.SendCreated(c, "Room created successfully", fiber.Map{
		"room": room,
	})
}

// GetRoom handles GET /api/rooms/:id
func (h *RoomHandler) GetRoom(c *fiber.Ctx) error {
	roomID := c.Params("id")
	userID := c.Locals("userID").(uuid.UUID).String()

	room, err := h.roomService.GetRoom(userID, roomID)
	if err != nil {
		if err == utils.ErrRoomNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Room not found", err)
		}
		if err == utils.ErrUnauthorized {
			return utils.SendUnauthorized(c, "You don't have access to this room")
		}
		return utils.SendInternalError(c, "Failed to fetch room", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Room fetched successfully", fiber.Map{
		"room": room,
	})
}

// JoinRoom handles POST /api/rooms/join
func (h *RoomHandler) JoinRoom(c *fiber.Ctx) error {
	var req dto.JoinRoomRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendBadRequest(c, "Validation failed", err)
	}

	userID := c.Locals("userID").(uuid.UUID).String()

	room, err := h.roomService.JoinRoom(userID, &req)
	if err != nil {
		if err.Error() == "room not found with this code" {
			return utils.SendError(c, fiber.StatusNotFound, "Room not found", err)
		}
		if err.Error() == "room is full" {
			return utils.SendError(c, fiber.StatusConflict, "Room is full", err)
		}
		return utils.SendInternalError(c, "Failed to join room", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Joined room successfully", fiber.Map{
		"room": room,
	})
}

// LeaveRoom handles POST /api/rooms/:id/leave
func (h *RoomHandler) LeaveRoom(c *fiber.Ctx) error {
	roomID := c.Params("id")
	userID := c.Locals("userID").(uuid.UUID).String()

	if err := h.roomService.LeaveRoom(userID, roomID); err != nil {
		if err.Error() == "you are not a participant in this room" {
			return utils.SendError(c, fiber.StatusForbidden, "You are not a participant", err)
		}
		return utils.SendInternalError(c, "Failed to leave room", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Left room successfully", nil)
}

// GetUserRooms handles GET /api/rooms
func (h *RoomHandler) GetUserRooms(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID).String()

	rooms, err := h.roomService.GetUserRooms(userID)
	if err != nil {
		return utils.SendInternalError(c, "Failed to fetch rooms", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Rooms fetched successfully", fiber.Map{
		"rooms": rooms,
	})
}

// DeleteRoom handles DELETE /api/rooms/:id
func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
	roomID := c.Params("id")
	userID := c.Locals("userID").(uuid.UUID).String()

	if err := h.roomService.DeleteRoom(userID, roomID); err != nil {
		if err == utils.ErrRoomNotFound {
			return utils.SendError(c, fiber.StatusNotFound, "Room not found", err)
		}
		if err == utils.ErrUnauthorized {
			return utils.SendUnauthorized(c, "Only room creator can delete")
		}
		return utils.SendInternalError(c, "Failed to delete room", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Room deleted successfully", nil)
}

// GetCodeSession handles GET /api/rooms/:id/code
func (h *RoomHandler) GetCodeSession(c *fiber.Ctx) error {
	roomID := c.Params("id")
	userID := c.Locals("userID").(uuid.UUID).String()

	session, err := h.roomService.GetCodeSession(userID, roomID)
	if err != nil {
		if err == utils.ErrUnauthorized {
			return utils.SendUnauthorized(c, "You don't have access to this room")
		}
		if err.Error() == "no active code session" {
			return utils.SendError(c, fiber.StatusNotFound, "No active code session", err)
		}
		return utils.SendInternalError(c, "Failed to fetch code session", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Code session fetched successfully", fiber.Map{
		"session": session,
	})
}

// UpdateCodeSession handles PUT /api/rooms/:id/code
func (h *RoomHandler) UpdateCodeSession(c *fiber.Ctx) error {
	roomID := c.Params("id")
	userID := c.Locals("userID").(uuid.UUID).String()

	var req dto.UpdateCodeSessionRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequest(c, "Invalid request payload", err)
	}

	session, err := h.roomService.UpdateCodeSession(userID, roomID, &req)
	if err != nil {
		if err == utils.ErrUnauthorized {
			return utils.SendUnauthorized(c, "You don't have access to this room")
		}
		return utils.SendInternalError(c, "Failed to update code session", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Code session updated successfully", fiber.Map{
		"session": session,
	})
}
