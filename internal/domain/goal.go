package domain

import "time"

type Goal struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Target      float64   `json:"target"`
	Initial     float64   `json:"initial"`
	Current     float64   `json:"current"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
