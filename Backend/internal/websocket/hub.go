package websocket

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients by room
	Rooms map[uuid.UUID]map[*Client]bool

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Broadcast messages to clients in a room
	Broadcast chan *Message
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[uuid.UUID]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient adds a client to a room
func (h *Hub) registerClient(client *Client) {
	// Create room if doesn't exist
	if _, exists := h.Rooms[client.RoomID]; !exists {
		h.Rooms[client.RoomID] = make(map[*Client]bool)
	}

	h.Rooms[client.RoomID][client] = true

	log.Printf("Client registered: User %s joined room %s", client.Username, client.RoomID)

	// Notify other users in room
	h.notifyUserJoined(client)
}

// unregisterClient removes a client from a room
func (h *Hub) unregisterClient(client *Client) {
	if room, exists := h.Rooms[client.RoomID]; exists {
		if _, ok := room[client]; ok {
			delete(room, client)
			close(client.Send)

			// Remove room if empty
			if len(room) == 0 {
				delete(h.Rooms, client.RoomID)
			}

			log.Printf("Client unregistered: User %s left room %s", client.Username, client.RoomID)

			// Notify other users
			h.notifyUserLeft(client)
		}
	}
}

// broadcastMessage sends a message to all clients in the same room
func (h *Hub) broadcastMessage(message *Message) {
	room, exists := h.Rooms[message.RoomID]
	if !exists {
		return
	}

	for client := range room {
		select {
		case client.Send <- message:
		default:
			// Client's send buffer is full, disconnect
			close(client.Send)
			delete(room, client)
		}
	}
}

// notifyUserJoined sends user joined notification
func (h *Hub) notifyUserJoined(newClient *Client) {
	room, exists := h.Rooms[newClient.RoomID]
	if !exists {
		return
	}

	// Send user list to new client
	users := h.getRoomUsers(newClient.RoomID)
	userListData, _ := json.Marshal(users)
	userListMsg := &Message{
		Type:   MessageTypeUserList,
		RoomID: newClient.RoomID,
		Data:   userListData,
	}
	newClient.Send <- userListMsg

	// Notify others about new user
	userInfo := UserInfo{
		UserID:   newClient.UserID,
		Username: newClient.Username,
		Color:    newClient.Color,
		IsOnline: true,
	}
	userData, _ := json.Marshal(userInfo)
	joinMsg := &Message{
		Type:     MessageTypeUserJoined,
		RoomID:   newClient.RoomID,
		UserID:   newClient.UserID,
		Username: newClient.Username,
		Data:     userData,
	}

	for client := range room {
		if client.ID != newClient.ID {
			client.Send <- joinMsg
		}
	}
}

// notifyUserLeft sends user left notification
func (h *Hub) notifyUserLeft(leftClient *Client) {
	room, exists := h.Rooms[leftClient.RoomID]
	if !exists {
		return
	}

	userInfo := UserInfo{
		UserID:   leftClient.UserID,
		Username: leftClient.Username,
		IsOnline: false,
	}
	userData, _ := json.Marshal(userInfo)
	leaveMsg := &Message{
		Type:     MessageTypeUserLeft,
		RoomID:   leftClient.RoomID,
		UserID:   leftClient.UserID,
		Username: leftClient.Username,
		Data:     userData,
	}

	for client := range room {
		client.Send <- leaveMsg
	}
}

// getRoomUsers returns list of users in a room
func (h *Hub) getRoomUsers(roomID uuid.UUID) []UserInfo {
	room, exists := h.Rooms[roomID]
	if !exists {
		return []UserInfo{}
	}

	users := make([]UserInfo, 0, len(room))
	for client := range room {
		users = append(users, UserInfo{
			UserID:   client.UserID,
			Username: client.Username,
			Color:    client.Color,
			IsOnline: true,
		})
	}

	return users
}
