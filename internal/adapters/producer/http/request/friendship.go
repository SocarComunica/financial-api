package request

import "time"

type SendFriendRequest struct {
	FriendID uint `json:"friend_id" validate:"required"`
}

type UpdateFriendshipStatus struct {
	Status string `json:"status" validate:"required,oneof=accepted rejected"`
}

// UserOutputDTO defines the standard output structure for user information.
type UserOutputDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type FriendStats struct {
	UserID          uint    `json:"user_id"`
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	TotalGoals      int     `json:"total_goals"`
	CompletedGoals  int     `json:"completed_goals"`
	AverageProgress float64 `json:"average_progress"`
}

// FriendshipOutputDTO defines the standard output for friendship details.
type FriendshipOutputDTO struct {
	ID        uint           `json:"id"`
	User1     *UserOutputDTO `json:"user1,omitempty"`
	User2     *UserOutputDTO `json:"user2,omitempty"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// UserAndFriendsStatsResponse is the structure for the combined user and friends stats response
type UserAndFriendsStatsResponse struct {
	UserStats    *FriendStats   `json:"user_stats"`
	FriendsStats []*FriendStats `json:"friends_stats"`
}
