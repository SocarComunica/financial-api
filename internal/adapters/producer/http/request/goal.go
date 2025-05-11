package request

// CreateGoal is the request body for creating a new goal
type CreateGoal struct {
	UserID      uint    `json:"user_id" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	Target      float64 `json:"target" validate:"required"`
	Initial     float64 `json:"initial"`
}

// UpdateGoal is the request body for updating an existing goal
type UpdateGoal struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Target      *float64 `json:"target,omitempty"`
	Current     *float64 `json:"current,omitempty"`
}
