package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Accounts []Account
	Goals    []Goal
	// Relaciones de amistad
	FriendshipsAsUser1 []Friendship `gorm:"foreignKey:User1ID"`
	FriendshipsAsUser2 []Friendship `gorm:"foreignKey:User2ID"`
}

type Friendship struct {
	gorm.Model
	User1ID uint
	User2ID uint
	Status  FriendshipStatus `gorm:"type:varchar(20);not null"`
	User1   User             `gorm:"foreignKey:User1ID"`
	User2   User             `gorm:"foreignKey:User2ID"`
}

type FriendshipStatus string

const (
	FriendshipStatusPending  FriendshipStatus = "pending"
	FriendshipStatusAccepted FriendshipStatus = "accepted"
	FriendshipStatusRejected FriendshipStatus = "rejected"
)
