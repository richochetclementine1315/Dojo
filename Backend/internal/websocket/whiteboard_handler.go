package websocket

// WhiteboardHandler handles realtime whiteboard collaboration
type WhiteboardHandler struct {
	Hub *Hub
}

// NewWhiteboardHandler creates a new WhiteboardHandler
func NewWhiteboardHandler(hub *Hub) *WhiteboardHandler {
	return &WhiteboardHandler{
		Hub: hub}
}
