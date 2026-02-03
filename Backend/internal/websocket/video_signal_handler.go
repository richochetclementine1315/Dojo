package websocket

// VideoSignalHandler handles WebRTC signaling for video calls
type VideoSignalHandler struct {
	Hub *Hub
}

// NewVideoSignalHandler creates a new VideoSignalHandler
func NewVideoSignalHandler(hub *Hub) *VideoSignalHandler {
	return &VideoSignalHandler{
		Hub: hub,
	}
}
