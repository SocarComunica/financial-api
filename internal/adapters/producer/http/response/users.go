package response

import "github.com/socarcomunica/financial-api/internal/domain"

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// MapUserToResponse creates a user response from a user domain model
func MapUserToResponse(user *domain.User) *User {
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
