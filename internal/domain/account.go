package domain

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Name           string  `json:"name" gorm:"not null"`
	Type           string  `json:"type" gorm:"not null"`
	Balance        float64 `json:"balance" gorm:"not null;default:0"`
	InitialBalance float64 `json:"initial_balance" gorm:"not null;default:0"`
	UserID         uint    `json:"user_id" gorm:"not null"`
	User           User    `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
