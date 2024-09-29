package domain

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}
