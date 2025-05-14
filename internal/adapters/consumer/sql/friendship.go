package sql

import (
	"errors"

	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
	"gorm.io/gorm"
)

const (
	CreateFriendshipError           = "CreateFriendshipError DB Client: "
	GetFriendshipError              = "GetFriendshipError DB Client: "
	UpdateFriendshipError           = "UpdateFriendshipError DB Client: "
	GetFriendsError                 = "GetFriendsError DB Client: "
	GetPendingFriendRequestsDBError = "GetPendingFriendRequestsError DB Client: "
	GetUserByIDDBError              = "GetUserByIDError DB Client: "
)

func (c *client) CreateFriendship(user1ID, user2ID uint) (*domain.Friendship, error) {
	// Verificar que ambos usuarios existan
	var user1, user2 domain.User
	if result := c.DB.First(&user1, user1ID); result.Error != nil {
		return nil, errors.New(CreateFriendshipError + "user1 not found")
	}
	if result := c.DB.First(&user2, user2ID); result.Error != nil {
		return nil, errors.New(CreateFriendshipError + "user2 not found")
	}

	// Verificar que no exista ya una amistad
	var existingFriendship domain.Friendship
	result := c.DB.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID).First(&existingFriendship)
	if result.Error == nil {
		return nil, errors.New(CreateFriendshipError + "friendship already exists")
	}

	friendship := &domain.Friendship{
		User1ID: user1ID,
		User2ID: user2ID,
		Status:  domain.FriendshipStatusPending,
	}

	if result := c.DB.Create(friendship); result.Error != nil {
		return nil, errors.New(CreateFriendshipError + result.Error.Error())
	}

	return friendship, nil
}

func (c *client) UpdateFriendshipStatus(id uint, status domain.FriendshipStatus) error {
	result := c.DB.Model(&domain.Friendship{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return errors.New(UpdateFriendshipError + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return errors.New(UpdateFriendshipError + "friendship not found")
	}
	return nil
}

func (c *client) GetFriends(userID uint) ([]*domain.User, error) {
	var friendships []*domain.Friendship
	result := c.DB.Where("(user1_id = ? OR user2_id = ?) AND status = ?",
		userID, userID, domain.FriendshipStatusAccepted).Find(&friendships)
	if result.Error != nil {
		return nil, errors.New(GetFriendsError + result.Error.Error())
	}

	var friends []*domain.User
	for _, f := range friendships {
		var friendID uint
		if f.User1ID == userID {
			friendID = f.User2ID
		} else {
			friendID = f.User1ID
		}

		var friend domain.User
		if result := c.DB.First(&friend, friendID); result.Error != nil {
			return nil, errors.New(GetFriendsError + "friend not found")
		}
		friends = append(friends, &friend)
	}

	return friends, nil
}

func (c *client) GetFriendStats(userID uint) (*request.FriendStats, error) {
	var user domain.User
	if result := c.DB.Preload("Goals").First(&user, userID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New(GetFriendsError + "user not found")
		}
		return nil, errors.New(GetFriendsError + "DB error fetching user: " + result.Error.Error())
	}

	stats := &request.FriendStats{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	totalGoals := len(user.Goals)
	stats.TotalGoals = totalGoals

	if totalGoals == 0 {
		return stats, nil
	}

	var totalProgress float64
	for _, goal := range user.Goals {
		if goal.Current >= goal.Target {
			stats.CompletedGoals++
		}
		progress := float64(goal.Current) / float64(goal.Target) * 100
		totalProgress += progress
	}

	stats.AverageProgress = totalProgress / float64(totalGoals)
	return stats, nil
}

func (c *client) GetPendingFriendRequests(userID uint) ([]*domain.Friendship, error) {
	var pendingRequests []*domain.Friendship
	// Find requests where the current user is the recipient (User2ID) and status is Pending.
	// It might be useful to Preload User1 (the sender) to get their details, e.g., Username.
	// Assuming Friendship has a User1 field/association for the sender.
	// If not, you'd just get User1ID.
	result := c.DB.Joins("User1").Where("friendships.user2_id = ? AND friendships.status = ?", userID, domain.FriendshipStatusPending).Find(&pendingRequests)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []*domain.Friendship{}, nil // No pending requests found, return empty slice
		}
		return nil, errors.New(GetPendingFriendRequestsDBError + result.Error.Error())
	}

	return pendingRequests, nil
}

// GetUserByID retrieves a user by their ID.
func (c *client) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	if result := c.DB.First(&user, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New(GetUserByIDDBError + "user not found")
		}
		return nil, errors.New(GetUserByIDDBError + result.Error.Error())
	}
	return &user, nil
}
