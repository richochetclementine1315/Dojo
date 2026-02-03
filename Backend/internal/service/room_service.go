package service

import (
	"dojo/internal/dto"
	"dojo/internal/models"
	"dojo/internal/repository"
	"dojo/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomService struct {
	roomRepo *repository.RoomRepository
	userRepo *repository.UserRepository
}

func NewRoomService(roomRepo *repository.RoomRepository, userRepo *repository.UserRepository) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
		userRepo: userRepo,
	}
}

// CreateRoom creates a new collaborative coding room
func (s *RoomService) CreateRoom(userID string, req *dto.CreateRoomRequest) (*dto.RoomResponse, error) {
	// Parse user ID
	creatorID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Generate unique room code
	roomCode, err := s.roomRepo.GenerateUniqueRoomCode()
	if err != nil {
		return nil, errors.New("failed to generate room code")
	}

	// Set default max participants if not provided
	maxParticipants := req.MaxParticipants
	if maxParticipants == 0 {
		maxParticipants = 4
	}

	// Create room
	room := &models.Room{
		Name:            req.Name,
		RoomCode:        roomCode,
		CreatedBy:       &creatorID,
		MaxParticipants: maxParticipants,
		IsActive:        true,
	}

	if err := s.roomRepo.Create(room); err != nil {
		return nil, err
	}

	// Add creator as first participant
	participant := &models.RoomParticipant{
		RoomID:   room.ID,
		UserID:   creatorID,
		IsOnline: true,
	}

	if err := s.roomRepo.AddParticipant(participant); err != nil {
		return nil, err
	}

	// Reload room with relationships
	room, err = s.roomRepo.FindByID(room.ID.String())
	if err != nil {
		return nil, err
	}

	// Create initial code session
	codeSession := &models.CodeSession{
		RoomID:   room.ID,
		Language: "javascript",
		Code:     "// Start coding here...\n",
	}
	if err := s.roomRepo.CreateCodeSession(codeSession); err != nil {
		return nil, err
	}

	return s.mapRoomToResponse(room), nil
}

// GetRoom retrieves a room by ID
func (s *RoomService) GetRoom(userID, roomID string) (*dto.RoomResponse, error) {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRoomNotFound
		}
		return nil, err
	}

	// Check if user is participant or creator
	isParticipant, err := s.roomRepo.IsParticipant(roomID, userID)
	if err != nil {
		return nil, err
	}

	if !isParticipant && room.CreatedBy.String() != userID {
		return nil, utils.ErrUnauthorized
	}

	return s.mapRoomToResponse(room), nil
}

// JoinRoom allows a user to join a room using room code
func (s *RoomService) JoinRoom(userID string, req *dto.JoinRoomRequest) (*dto.RoomResponse, error) {
	// Find room by code
	room, err := s.roomRepo.FindByRoomCode(req.RoomCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("room not found with this code")
		}
		return nil, err
	}

	// Check if user is already a participant
	isParticipant, err := s.roomRepo.IsParticipant(room.ID.String(), userID)
	if err != nil {
		return nil, err
	}

	if isParticipant {
		return s.mapRoomToResponse(room), nil
	}

	// Check if room is full
	activeCount, err := s.roomRepo.GetActiveParticipantCount(room.ID.String())
	if err != nil {
		return nil, err
	}

	if int(activeCount) >= room.MaxParticipants {
		return nil, errors.New("room is full")
	}

	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Add user as participant
	participant := &models.RoomParticipant{
		RoomID:   room.ID,
		UserID:   userUUID,
		IsOnline: true,
	}

	if err := s.roomRepo.AddParticipant(participant); err != nil {
		return nil, err
	}

	// Reload room with updated participants
	room, err = s.roomRepo.FindByID(room.ID.String())
	if err != nil {
		return nil, err
	}

	return s.mapRoomToResponse(room), nil
}

// LeaveRoom allows a user to leave a room
func (s *RoomService) LeaveRoom(userID, roomID string) error {
	// Check if user is participant
	isParticipant, err := s.roomRepo.IsParticipant(roomID, userID)
	if err != nil {
		return err
	}

	if !isParticipant {
		return errors.New("you are not a participant in this room")
	}

	return s.roomRepo.RemoveParticipant(roomID, userID)
}

// GetUserRooms retrieves all rooms where user is creator or participant
func (s *RoomService) GetUserRooms(userID string) ([]dto.RoomResponse, error) {
	rooms, err := s.roomRepo.FindUserRooms(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.RoomResponse, len(rooms))
	for i, room := range rooms {
		responses[i] = *s.mapRoomToResponse(&room)
	}

	return responses, nil
}

// DeleteRoom deletes a room (only creator can delete)
func (s *RoomService) DeleteRoom(userID, roomID string) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrRoomNotFound
		}
		return err
	}

	// Check if user is creator
	if room.CreatedBy == nil || room.CreatedBy.String() != userID {
		return utils.ErrUnauthorized
	}

	return s.roomRepo.Delete(roomID)
}

// GetCodeSession retrieves the active code session for a room
func (s *RoomService) GetCodeSession(userID, roomID string) (*dto.CodeSessionResponse, error) {
	// Check if user is participant
	isParticipant, err := s.roomRepo.IsParticipant(roomID, userID)
	if err != nil {
		return nil, err
	}

	if !isParticipant {
		return nil, utils.ErrUnauthorized
	}

	session, err := s.roomRepo.GetCodeSession(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active code session")
		}
		return nil, err
	}

	return s.mapCodeSessionToResponse(session), nil
}

// UpdateCodeSession updates the code session
func (s *RoomService) UpdateCodeSession(userID, roomID string, req *dto.UpdateCodeSessionRequest) (*dto.CodeSessionResponse, error) {
	// Check if user is participant
	isParticipant, err := s.roomRepo.IsParticipant(roomID, userID)
	if err != nil {
		return nil, err
	}

	if !isParticipant {
		return nil, utils.ErrUnauthorized
	}

	session, err := s.roomRepo.GetCodeSession(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new session
			roomUUID, _ := uuid.Parse(roomID)
			session = &models.CodeSession{
				RoomID:    roomUUID,
				Language:  req.Language,
				Code:      req.Code,
				ProblemID: req.ProblemID,
			}
			if err := s.roomRepo.CreateCodeSession(session); err != nil {
				return nil, err
			}
			return s.mapCodeSessionToResponse(session), nil
		}
		return nil, err
	}

	// Update session
	if req.Language != "" {
		session.Language = req.Language
	}
	if req.Code != "" {
		session.Code = req.Code
	}
	if req.ProblemID != nil {
		session.ProblemID = req.ProblemID
	}

	if err := s.roomRepo.UpdateCodeSession(session); err != nil {
		return nil, err
	}

	return s.mapCodeSessionToResponse(session), nil
}

// mapRoomToResponse converts Room model to RoomResponse DTO
func (s *RoomService) mapRoomToResponse(room *models.Room) *dto.RoomResponse {
	response := &dto.RoomResponse{
		ID:              room.ID,
		Name:            room.Name,
		RoomCode:        room.RoomCode,
		MaxParticipants: room.MaxParticipants,
		IsActive:        room.IsActive,
		CreatedAt:       room.CreatedAt,
	}

	if room.Creator != nil {
		response.Creator = &dto.UserResponse{
			ID:         room.Creator.ID,
			Username:   room.Creator.Username,
			Email:      room.Creator.Email,
			AvatarURL:  room.Creator.AvatarURL,
			IsVerified: room.Creator.IsVerified,
			CreatedAt:  room.Creator.CreatedAt,
		}
	}

	participants := make([]dto.ParticipantResponse, len(room.Participants))
	for i, p := range room.Participants {
		if p.LeftAt == nil {
			participants[i] = dto.ParticipantResponse{
				ID:       p.ID,
				IsOnline: p.IsOnline,
				JoinedAt: p.JoinedAt,
				User: dto.UserResponse{
					ID:         p.User.ID,
					Username:   p.User.Username,
					Email:      p.User.Email,
					AvatarURL:  p.User.AvatarURL,
					IsVerified: p.User.IsVerified,
					CreatedAt:  p.User.CreatedAt,
				},
			}
		}
	}
	response.Participants = participants

	return response
}

// mapCodeSessionToResponse converts CodeSession model to CodeSessionResponse DTO
func (s *RoomService) mapCodeSessionToResponse(session *models.CodeSession) *dto.CodeSessionResponse {
	response := &dto.CodeSessionResponse{
		ID:        session.ID,
		RoomID:    session.RoomID,
		Language:  session.Language,
		Code:      session.Code,
		UpdatedAt: session.UpdatedAt,
	}

	if session.Problem != nil {
		response.Problem = &dto.ProblemResponse{
			ID:                session.Problem.ID,
			Platform:          session.Problem.Platform,
			PlatformProblemID: session.Problem.PlatformProblemID,
			Title:             session.Problem.Title,
			Slug:              session.Problem.Slug,
			Difficulty:        session.Problem.Difficulty,
			Tags:              session.Problem.Tags,
			AcceptanceRate:    session.Problem.AcceptanceRate,
			ProblemURL:        session.Problem.ProblemURL,
			Description:       session.Problem.Description,
			Constraints:       session.Problem.Constraints,
			Examples:          session.Problem.Examples,
			Hints:             session.Problem.Hints,
			CreatedAt:         session.Problem.CreatedAt,
		}
	}

	return response
}
