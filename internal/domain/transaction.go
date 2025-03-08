package domain

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Amount        float64  `json:"amount"`
	Description   string   `json:"description"`
	Tags          []Tag    `json:"tags" gorm:"many2many:transaction_tags;"`
	Type          string   `json:"type"`
	OriginID      uint     `json:"origin_id"`
	Origin        Account  `json:"origin"`
	DestinationID *uint    `json:"destination_id"`
	Destination   *Account `json:"destination"`
}
