package utils

import "errors"

// Custom error types for better error handling
var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrTokenExpired       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("invalid token")

	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailTaken        = errors.New("email already in use")
	ErrUsernameTaken     = errors.New("username already in use")

	// Friend errors
	ErrFriendRequestNotFound      = errors.New("friend request not found")
	ErrFriendRequestAlreadyExists = errors.New("friend request already sent")
	ErrAlreadyFriends             = errors.New("users are already friends")
	ErrCannotSendToSelf           = errors.New("cannot send friend request to yourself")
	ErrUserBlocked                = errors.New("user is blocked")

	// Problem errors
	ErrProblemNotFound = errors.New("problem not found")
	ErrNoteNotFound    = errors.New("note not found")

	// Sheet errors
	ErrSheetNotFound         = errors.New("sheet not found")
	ErrSheetAccessDenied     = errors.New("access denied to this sheet")
	ErrProblemAlreadyInSheet = errors.New("problem already exists in this sheet")

	// Contest errors
	ErrContestNotFound = errors.New("contest not found")

	// Room errors
	ErrRoomNotFound  = errors.New("room not found")
	ErrRoomFull      = errors.New("room is full")
	ErrRoomInactive  = errors.New("room is inactive")
	ErrNotInRoom     = errors.New("user is not in this room")
	ErrAlreadyInRoom = errors.New("user is already in this room")

	// General errors
	ErrInvalidInput   = errors.New("invalid input")
	ErrDatabaseError  = errors.New("database error occurred")
	ErrInternalServer = errors.New("internal server error")
)
