package websocket

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	// connection messages
	MessageTypeJoin  MessageType = "join"
	MessageTypeLeave MessageType = "leave"
	MessageTypePing  MessageType = "ping" //You alive? Or should I kick u out :)
	MessageTypePong  MessageType = "pong" //Yaa  bro doing good... ;)

	// Code collaboration messages
	MessageTypeCodeUpdate     MessageType = "code_edit"
	MessageTypeCursorMove     MessageType = "cursor_move"
	MessageTypeCodeSelection  MessageType = "code_selection"
	MessageTypeLanguageChange MessageType = "language_change"

	// Chat messages
	MessageTypeChat MessageType = "chat"

	// WhiteBoard messages
	MessageTypeWhiteBoardDraw  MessageType = "whiteboard_draw"
	MessageTypeWhiteBoardClear MessageType = "whiteboard_clear"
	MessageTypeWhiteBoardUndo  MessageType = "whiteboard_undo"

	// Video/Audio WebRTC signaling messages
	MessageTypeRTCOffer     MessageType = "rtc_offer"
	MessageTypeRTCAnswer    MessageType = "rtc_answer"
	MessageTypeRTCCandidate MessageType = "rtc_candidate"

	// Room State messages
	MessageTypeUserJoined MessageType = "user_joined"
	MessageTypeUserLeft   MessageType = "user_left"
	MessageTypeUserList   MessageType = "user_list"

	// Error messages
	MessageTypeError MessageType = "error"
)

// message struct represents a websocket message
type Message struct {
	Type      MessageType     `json:"type"`
	RoomID    uuid.UUID       `json: "room_id"`
	UserID    uuid.UUID       `json: "user_id"`
	Username  string          `json: "username"`
	Data      json.RawMessage `json: "data"`
	Timestamp time.Time       `json: "timestamp"`
}

// CodeUpdateData represents the code editor update data
type CodeUpdateData struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Version  int    `json:"version"`
}
type CursorMoveData struct {
	Line   int    `json:"line"`
	Column int    `json:"column"`
	Color  string `json:"color"` //users cursor color..everybody will have different different colour...
}

// CodeSelection
type CodeSelectionData struct {
	StartLine   int `json:"start_line"`
	StartColumn int `json:"start_column"`
	EndLine     int `json:"end_line"`
	EndColumn   int `json:"end_column"`
}

// ChatData
type ChatData struct {
	Message string `json:"message"`
}

// WhiteboardDrawData represents whiteboard drawing data
type WhiteboardDrawData struct {
	Tool   string  `json:"tool"` // e.g., "pen", "eraser"
	Color  string  `json:"color"`
	Width  int     `json:"width"`
	Points []Point `json:"points"`
	Action string  `json:"action"` // e.g.,"start", "end","move"
}

// Point Represents THE CORDINATES(very important)...
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// RTCSignalData represents WebRTC signaling data
type RTCSignalData struct {
	TargetUserID uuid.UUID       `json:"target_user_id"`
	Signal       json.RawMessage `json:"signal"`
}
type UserInfo struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Color    string    `json:"color"` //Assigned cursor color
	IsOnline bool      `json:"is_online"`
}

// Errordata
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
