package request

type CreateAccount struct {
	Name    string  `json:"name" validate:"required"`
	Type    string  `json:"type" validate:"required"`
	Balance float64 `json:"balance" validate:"required"`
}
