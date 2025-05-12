package friendship

import (
	"errors"

	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateFriendshipServiceError = "CreateFriendship Service Error: "
	UpdateFriendshipServiceError = "UpdateFriendship Service Error: "
	GetFriendsServiceError       = "GetFriends Service Error: "
)

type FriendshipDatabase interface {
	CreateFriendship(user1ID, user2ID uint) (*domain.Friendship, error)
	UpdateFriendshipStatus(id uint, status domain.FriendshipStatus) error
	GetFriends(userID uint) ([]*domain.User, error)
	GetFriendStats(userID uint) (*request.FriendStats, error)
}

type Service struct {
	Database FriendshipDatabase
}

func NewFriendshipService(database FriendshipDatabase) *Service {
	return &Service{
		Database: database,
	}
}

func (s *Service) SendFriendRequest(userID, friendID uint) (*domain.Friendship, error) {
	if userID == friendID {
		return nil, errors.New(CreateFriendshipServiceError + "cannot send friend request to yourself")
	}

	friendship, err := s.Database.CreateFriendship(userID, friendID)
	if err != nil {
		return nil, errors.New(CreateFriendshipServiceError + err.Error())
	}

	return friendship, nil
}

func (s *Service) UpdateFriendshipStatus(id uint, status string) error {
	friendshipStatus := domain.FriendshipStatus(status)
	if friendshipStatus != domain.FriendshipStatusAccepted && friendshipStatus != domain.FriendshipStatusRejected {
		return errors.New(UpdateFriendshipServiceError + "invalid status")
	}

	err := s.Database.UpdateFriendshipStatus(id, friendshipStatus)
	if err != nil {
		return errors.New(UpdateFriendshipServiceError + err.Error())
	}

	return nil
}

func (s *Service) GetFriends(userID uint) ([]*domain.User, error) {
	friends, err := s.Database.GetFriends(userID)
	if err != nil {
		return nil, errors.New(GetFriendsServiceError + err.Error())
	}
	return friends, nil
}

func (s *Service) GetFriendStats(userID uint) (*request.FriendStats, error) {
	stats, err := s.Database.GetFriendStats(userID)
	if err != nil {
		return nil, errors.New(GetFriendsServiceError + err.Error())
	}
	return stats, nil
}
