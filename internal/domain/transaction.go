package domain

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Tags        []Tag   `json:"tags" gorm:"many2many:transaction_tags;"`
	Type        string  `json:"type"`
	Origin      uint    `json:"origin"`
	Destination *uint   `json:"destination"`
}
