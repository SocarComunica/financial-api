package request

type CreateTransaction struct {
	Amount        float64                `json:"amount" validate:"required,positive"`
	Description   string                 `json:"description"`
	Tags          []CreateTransactionTag `json:"tags"`
	Type          TransactionType        `json:"type"`
	OriginID      uint                   `json:"origin_id"`
	DestinationID *uint                  `json:"destination_id"`
}

type TransactionType struct {
	Name string `json:"name" validate:"required,transaction_type"`
}

type CreateTransactionTag struct {
	Name string `json:"name" validate:"required"`
}
