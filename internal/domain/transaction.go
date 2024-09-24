package domain

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Tags        string  `json:"tags"`
	Type        string  `json:"type"`
	Origin      uint    `json:"origin"`
	Destination uint    `json:"destination"`
}
