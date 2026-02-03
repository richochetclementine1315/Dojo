package repository

import (
	"dojo/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// Create creates a new room
func (r *RoomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

// FindByID retrieves a room by ID with all relationships
func (r *RoomRepository) FindByID(id string) (*models.Room, error) {
	var room models.Room
	err := r.db.Preload("Creator").
		Preload("Participants.User").
		Where("id = ?", id).
		First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// FindByRoomCode retrieves a room by room code
func (r *RoomRepository) FindByRoomCode(roomCode string) (*models.Room, error) {
	var room models.Room
	err := r.db.Preload("Creator").
		Preload("Participants.User").
		Where("room_code = ? AND is_active = ?", roomCode, true).
		First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// FindUserRooms retrieves all rooms where user is creator or participant
func (r *RoomRepository) FindUserRooms(userID string) ([]models.Room, error) {
	var rooms []models.Room

	err := r.db.Preload("Creator").
		Preload("Participants.User").
		Where("created_by = ? AND is_active = ?", userID, true).
		Or("id IN (SELECT room_id FROM room_participants WHERE user_id = ? AND left_at IS NULL)", userID).
		Find(&rooms).Error

	return rooms, err
}

// Delete soft deletes a room by setting is_active to false
func (r *RoomRepository) Delete(id string) error {
	return r.db.Model(&models.Room{}).Where("id = ?", id).Update("is_active", false).Error
}

// AddParticipant adds a participant to a room
func (r *RoomRepository) AddParticipant(participant *models.RoomParticipant) error {
	return r.db.Create(participant).Error
}

// RemoveParticipant marks a participant as left
func (r *RoomRepository) RemoveParticipant(roomID, userID string) error {
	now := time.Now()
	return r.db.Model(&models.RoomParticipant{}).
		Where("room_id = ? AND user_id = ? AND left_at IS NULL", roomID, userID).
		Updates(map[string]interface{}{
			"left_at":   &now,
			"is_online": false,
		}).Error
}

// GetActiveParticipantCount gets the number of active participants in a room
func (r *RoomRepository) GetActiveParticipantCount(roomID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.RoomParticipant{}).
		Where("room_id = ? AND left_at IS NULL", roomID).
		Count(&count).Error
	return count, err
}

// IsParticipant checks if a user is a participant in a room
func (r *RoomRepository) IsParticipant(roomID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.RoomParticipant{}).
		Where("room_id = ? AND user_id = ? AND left_at IS NULL", roomID, userID).
		Count(&count).Error
	return count > 0, err
}

// UpdateParticipantStatus updates participant online status
func (r *RoomRepository) UpdateParticipantStatus(roomID, userID string, isOnline bool) error {
	return r.db.Model(&models.RoomParticipant{}).
		Where("room_id = ? AND user_id = ? AND left_at IS NULL", roomID, userID).
		Update("is_online", isOnline).Error
}

// CreateCodeSession creates a new code session for a room
func (r *RoomRepository) CreateCodeSession(session *models.CodeSession) error {
	return r.db.Create(session).Error
}

// GetCodeSession retrieves the active code session for a room
func (r *RoomRepository) GetCodeSession(roomID string) (*models.CodeSession, error) {
	var session models.CodeSession
	err := r.db.Preload("Problem").
		Where("room_id = ?", roomID).
		Order("updated_at DESC").
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateCodeSession updates a code session
func (r *RoomRepository) UpdateCodeSession(session *models.CodeSession) error {
	return r.db.Save(session).Error
}

// GenerateUniqueRoomCode generates a unique 6-character room code
func (r *RoomRepository) GenerateUniqueRoomCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	for attempts := 0; attempts < 10; attempts++ {
		code := ""
		uid := uuid.New().String()
		for i := 0; i < codeLength; i++ {
			code += string(charset[uid[i]%byte(len(charset))])
		}

		// Check if code already exists
		var count int64
		if err := r.db.Model(&models.Room{}).Where("room_code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return code, nil
		}
	}

	return "", gorm.ErrInvalidData
}
