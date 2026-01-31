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

type SocialService struct {
	socialRepo *repository.SocialRepository
	userRepo   *repository.UserRepository
}

func NewSocialService(socialRepo *repository.SocialRepository, userRepo *repository.UserRepository) *SocialService {
	return &SocialService{
		socialRepo: socialRepo,
		userRepo:   userRepo,
	}
}

// SendFriendRequest sends a friend request
func (s *SocialService) SendFriendRequest(senderID string, req *dto.SendFriendRequestRequest) (*dto.FriendRequestResponse, error) {
	// Check if trying to send to self
	if senderID == req.ReceiverID.String() {
		return nil, utils.ErrCannotSendToSelf
	}

	// Check if receiver exists
	_, err := s.userRepo.FindByID(req.ReceiverID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}

	// Check if already friends
	areFriends, err := s.socialRepo.AreFriends(senderID, req.ReceiverID.String())
	if err != nil {
		return nil, err
	}
	if areFriends {
		return nil, utils.ErrAlreadyFriends
	}

	// Check if blocked
	isBlocked, err := s.socialRepo.IsBlocked(req.ReceiverID.String(), senderID)
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, utils.ErrUserBlocked
	}

	// Check if pending request already exists
	_, err = s.socialRepo.FindPendingRequest(senderID, req.ReceiverID.String())
	if err == nil {
		return nil, utils.ErrFriendRequestAlreadyExists
	}

	senderUUID, _ := uuid.Parse(senderID)
	friendRequest := &models.FriendRequest{
		SenderID:   senderUUID,
		ReceiverID: req.ReceiverID,
		Status:     "pending",
	}

	if err := s.socialRepo.CreateFriendRequest(friendRequest); err != nil {
		return nil, err
	}

	// Reload with user data
	friendRequest, err = s.socialRepo.FindFriendRequestByID(friendRequest.ID.String())
	if err != nil {
		return nil, err
	}

	return s.mapFriendRequestToResponse(friendRequest), nil
}

// GetReceivedRequests gets all pending friend requests received by a user
func (s *SocialService) GetReceivedRequests(userID string) ([]dto.FriendRequestResponse, error) {
	requests, err := s.socialRepo.GetReceivedRequests(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.FriendRequestResponse, len(requests))
	for i, req := range requests {
		responses[i] = *s.mapFriendRequestToResponse(&req)
	}

	return responses, nil
}

// GetSentRequests gets all pending friend requests sent by a user
func (s *SocialService) GetSentRequests(userID string) ([]dto.FriendRequestResponse, error) {
	requests, err := s.socialRepo.GetSentRequests(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.FriendRequestResponse, len(requests))
	for i, req := range requests {
		responses[i] = *s.mapFriendRequestToResponse(&req)
	}

	return responses, nil
}

// RespondToFriendRequest accepts or rejects a friend request
func (s *SocialService) RespondToFriendRequest(userID, requestID string, req *dto.UpdateFriendRequestRequest) error {
	friendRequest, err := s.socialRepo.FindFriendRequestByID(requestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrFriendRequestNotFound
		}
		return err
	}

	// Check if user is the receiver
	if friendRequest.ReceiverID.String() != userID {
		return utils.ErrUnauthorized
	}

	// Check if already processed
	if friendRequest.Status != "pending" {
		return errors.New("friend request already processed")
	}

	if req.Action == "accept" {
		// Create friendship
		if err := s.socialRepo.CreateFriend(friendRequest.SenderID.String(), friendRequest.ReceiverID.String()); err != nil {
			return err
		}

		friendRequest.Status = "accepted"
	} else {
		friendRequest.Status = "rejected"
	}

	return s.socialRepo.UpdateFriendRequest(friendRequest)
}

// CancelFriendRequest cancels a sent friend request
func (s *SocialService) CancelFriendRequest(userID, requestID string) error {
	friendRequest, err := s.socialRepo.FindFriendRequestByID(requestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrFriendRequestNotFound
		}
		return err
	}

	// Check if user is the sender
	if friendRequest.SenderID.String() != userID {
		return utils.ErrUnauthorized
	}

	return s.socialRepo.DeleteFriendRequest(requestID)
}

// GetFriends gets all friends of a user
func (s *SocialService) GetFriends(userID string) ([]dto.FriendResponse, error) {
	friends, err := s.socialRepo.GetFriends(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.FriendResponse, len(friends))
	for i, friend := range friends {
		responses[i] = *s.mapFriendToResponse(&friend)
	}

	return responses, nil
}

// RemoveFriend removes a friend
func (s *SocialService) RemoveFriend(userID, friendID string) error {
	// Check if they are friends
	areFriends, err := s.socialRepo.AreFriends(userID, friendID)
	if err != nil {
		return err
	}
	if !areFriends {
		return errors.New("not friends with this user")
	}

	return s.socialRepo.DeleteFriend(userID, friendID)
}

// BlockUser blocks a user
func (s *SocialService) BlockUser(userID string, req *dto.BlockUserRequest) error {
	// Check if trying to block self
	if userID == req.UserID.String() {
		return errors.New("cannot block yourself")
	}

	// Check if user exists
	_, err := s.userRepo.FindByID(req.UserID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrUserNotFound
		}
		return err
	}

	// Check if already blocked
	isBlocked, err := s.socialRepo.IsBlocked(userID, req.UserID.String())
	if err != nil {
		return err
	}
	if isBlocked {
		return errors.New("user already blocked")
	}

	// Remove friendship if exists
	areFriends, _ := s.socialRepo.AreFriends(userID, req.UserID.String())
	if areFriends {
		s.socialRepo.DeleteFriend(userID, req.UserID.String())
	}

	return s.socialRepo.CreateBlock(userID, req.UserID.String())
}

// UnblockUser unblocks a user
func (s *SocialService) UnblockUser(userID, blockedID string) error {
	// Check if blocked
	isBlocked, err := s.socialRepo.IsBlocked(userID, blockedID)
	if err != nil {
		return err
	}
	if !isBlocked {
		return errors.New("user not blocked")
	}

	return s.socialRepo.DeleteBlock(userID, blockedID)
}

// GetBlockedUsers gets all blocked users
func (s *SocialService) GetBlockedUsers(userID string) ([]dto.UserResponse, error) {
	blocks, err := s.socialRepo.GetBlockedUsers(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, len(blocks))
	for i, block := range blocks {
		responses[i] = *s.mapUserToResponse(&block.Blocked)
	}

	return responses, nil
}

// SearchUsers searches for users
func (s *SocialService) SearchUsers(query string, limit int) ([]dto.UserResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	users, err := s.socialRepo.SearchUsers(query, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *s.mapUserToResponse(&user)
	}

	return responses, nil
}

// Mapping functions

func (s *SocialService) mapFriendRequestToResponse(req *models.FriendRequest) *dto.FriendRequestResponse {
	return &dto.FriendRequestResponse{
		ID:        req.ID,
		Sender:    *s.mapUserToResponse(&req.Sender),
		Receiver:  *s.mapUserToResponse(&req.Receiver),
		Status:    req.Status,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}
}

func (s *SocialService) mapFriendToResponse(friend *models.Friend) *dto.FriendResponse {
	return &dto.FriendResponse{
		ID:        friend.ID,
		Friend:    *s.mapUserToResponse(&friend.Friend),
		IsOnline:  false, // TODO: Implement online status
		CreatedAt: friend.CreatedAt,
	}
}

func (s *SocialService) mapUserToResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		AvatarURL:  user.AvatarURL,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
	}
}
