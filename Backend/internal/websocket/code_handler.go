package websocket

// CodeHandler handles realtime code collaboration
type CodeHandler struct {
	Hub *Hub
}

// NewCodeHandler creates a new CodeHandler
func NewCodeHandler(hub *Hub) *CodeHandler {
	return &CodeHandler{
		Hub: hub,
	}
}
