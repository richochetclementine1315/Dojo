package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512KB
)

// Client represents a WebSocket client connection
type Client struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Username string
	RoomID   uuid.UUID
	Conn     *websocket.Conn
	Hub      *Hub
	Send     chan *Message
	Color    string // Assigned cursor color
}

// NewClient creates a new WebSocket client
func NewClient(userID uuid.UUID, username string, roomID uuid.UUID, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:       uuid.New(),
		UserID:   userID,
		Username: username,
		RoomID:   roomID,
		Conn:     conn,
		Hub:      hub,
		Send:     make(chan *Message, 256),
		Color:    assignColor(), // Assign random color for cursor
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Set metadata
		msg.UserID = c.UserID
		msg.Username = c.Username
		msg.RoomID = c.RoomID
		msg.Timestamp = time.Now()

		// Send to hub for broadcasting
		c.Hub.Broadcast <- &msg
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// assignColor returns a random color for cursor
func assignColor() string {
	colors := []string{
		"#FF6B6B", "#4ECDC4", "#45B7D1", "#FFA07A", "#98D8C8",
		"#F7DC6F", "#BB8FCE", "#85C1E2", "#F8B195", "#C06C84",
	}
	return colors[time.Now().UnixNano()%int64(len(colors))]
}
