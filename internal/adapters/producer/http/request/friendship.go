package request

type SendFriendRequest struct {
	FriendID uint `json:"friend_id" validate:"required"`
}

type UpdateFriendshipStatus struct {
	Status string `json:"status" validate:"required,oneof=accepted rejected"`
}

type FriendStats struct {
	UserID          uint    `json:"user_id"`
	Username        string  `json:"username"`
	TotalGoals      int     `json:"total_goals"`
	CompletedGoals  int     `json:"completed_goals"`
	AverageProgress float64 `json:"average_progress"`
}
