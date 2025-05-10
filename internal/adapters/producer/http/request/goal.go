package request

type CreateGoal struct {
	UserID      uint    `json:"user_id" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Target      float64 `json:"target" binding:"required"`
	Initial     float64 `json:"initial" binding:"required"`
}
