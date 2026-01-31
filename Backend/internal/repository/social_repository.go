package repository

import (
	"dojo/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SocialRepository struct {
	db *gorm.DB
}

func NewSocialRepository(db *gorm.DB) *SocialRepository {
	return &SocialRepository{db: db}
}

// Friend request stuffs
// CreateFriendRequest creates a new friend request
func (r *SocialRepository) CreateFriendRequest(request *models.FriendRequest) error {
	return r.db.Create(request).Error
}

// FindFriendRequestByID retrieves a friend request by ID
func (r *SocialRepository) FindFriendRequestByID(id string) (*models.FriendRequest, error) {
	var request models.FriendRequest
	err := r.db.Preload("Sender").Preload("Receiver").Where("id=?", id).First(&request).Error
	return &request, err
}

// FindPendingFriendRequests finds a list of pending friend requests for a user
func (r *SocialRepository) FindPendingRequest(senderID, receiverID string) (*models.FriendRequest, error) {
	var request models.FriendRequest
	err := r.db.Where("sender_id=? AND receiver_id=? AND status=?", senderID, receiverID, "pending").First(&request).Error
	return &request, err
}

// GetReceivedRequests gets all pending friend requests received by a user
func (r *SocialRepository) GetReceivedRequests(userID string) ([]models.FriendRequest, error) {
	var requests []models.FriendRequest
	err := r.db.Preload("Sender").Where("receiver_id = ? AND status = ?", userID, "pending").Order("created_at DESC").Find(&requests).Error
	return requests, err
}

// GetSentRequests gets all pending friend requests sent by a user
func (r *SocialRepository) GetSentRequests(userID string) ([]models.FriendRequest, error) {
	var requests []models.FriendRequest
	err := r.db.Preload("Receiver").Where("sender_id = ? AND status = ?", userID, "pending").Order("created_at DESC").Find(&requests).Error
	return requests, err
}

// UpdateFriendRequest updates a friend request
func (r *SocialRepository) UpdateFriendRequest(request *models.FriendRequest) error {
	return r.db.Save(request).Error
}

// DeleteFriendRequest deletes a friend request
func (r *SocialRepository) DeleteFriendRequest(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.FriendRequest{}).Error
}

// Friend Operations

// CreateFriend creates a friendship (bidirectional)
func (r *SocialRepository) CreateFriend(userID, friendID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create both directions
		friend1 := &models.Friend{UserID: parseUUID(userID), FriendID: parseUUID(friendID)}
		if err := tx.Create(friend1).Error; err != nil {
			return err
		}

		friend2 := &models.Friend{UserID: parseUUID(friendID), FriendID: parseUUID(userID)}
		if err := tx.Create(friend2).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetFriends gets all friends of a user
func (r *SocialRepository) GetFriends(userID string) ([]models.Friend, error) {
	var friends []models.Friend
	err := r.db.Preload("Friend").Where("user_id = ?", userID).Order("created_at DESC").Find(&friends).Error
	return friends, err
}

// AreFriends checks if two users are friends
func (r *SocialRepository) AreFriends(userID, friendID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Friend{}).Where("user_id = ? AND friend_id = ?", userID, friendID).Count(&count).Error
	return count > 0, err
}

// DeleteFriend removes a friendship (bidirectional)
func (r *SocialRepository) DeleteFriend(userID, friendID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete both directions
		if err := tx.Where("user_id = ? AND friend_id = ?", userID, friendID).Delete(&models.Friend{}).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ? AND friend_id = ?", friendID, userID).Delete(&models.Friend{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// Block Features

// CreateBlock creates a block relationship
func (r *SocialRepository) CreateBlock(blockerID, blockedID string) error {
	block := &models.BlockedUser{
		BlockerID: parseUUID(blockerID),
		BlockedID: parseUUID(blockedID),
	}
	return r.db.Create(block).Error
}

// IsBlocked checks if one user has blocked another
func (r *SocialRepository) IsBlocked(blockerID, blockedID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.BlockedUser{}).Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).Count(&count).Error
	return count > 0, err
}

// GetBlockedUsers gets all users blocked by a user
func (r *SocialRepository) GetBlockedUsers(blockerID string) ([]models.BlockedUser, error) {
	var blocks []models.BlockedUser
	err := r.db.Preload("Blocked").Where("blocker_id = ?", blockerID).Order("created_at DESC").Find(&blocks).Error
	return blocks, err
}

// DeleteBlock removes a block relationship
func (r *SocialRepository) DeleteBlock(blockerID, blockedID string) error {
	return r.db.Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).Delete(&models.BlockedUser{}).Error
}

// Search Users
func (r *SocialRepository) SearchUsers(query string, limit int) ([]models.User, error) {
	var users []models.User
	searchPattern := "%" + query + "%"
	err := r.db.Where("username ILIKE ? OR full_name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Find(&users).Error
	return users, err
}

// Helper function
func parseUUID(id string) uuid.UUID {
	uid, _ := uuid.Parse(id)
	return uid
}
