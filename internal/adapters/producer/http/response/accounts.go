package response

import "github.com/socarcomunica/financial-api/internal/domain"

type Account struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	Balance        float64 `json:"balance"`
	InitialBalance float64 `json:"initial_balance"`
	UserID         uint    `json:"user_id"`
}

// MapAccountToResponse creates an account response from an account domain model
func MapAccountToResponse(account *domain.Account) *Account {
	return &Account{
		ID:             account.ID,
		Name:           account.Name,
		Type:           account.Type,
		Balance:        account.Balance,
		InitialBalance: account.InitialBalance,
		UserID:         account.UserID,
	}
}
