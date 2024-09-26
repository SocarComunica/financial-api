package request

type CreateTransaction struct {
	Amount      float64                `json:"amount" validate:"required"`
	Description string                 `json:"description"`
	Tags        []CreateTransactionTag `json:"tags"`
	Type        string                 `json:"type"`
	Origin      uint                   `json:"origin"`
	Destination *uint                  `json:"destination"`
}

type CreateTransactionTag struct {
	Name string `json:"name" validate:"required"`
}
