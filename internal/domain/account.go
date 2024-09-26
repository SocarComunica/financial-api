package domain

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
}
