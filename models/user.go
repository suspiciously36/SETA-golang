package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	UserID       uint           `json:"userId" gorm:"primaryKey;autoIncrement"`
	Username     string         `json:"username" gorm:"unique;not null"`
	Email        string         `json:"email" gorm:"unique;not null"`
	Role         string         `json:"role" gorm:"not null"`
	PasswordHash string         `json:"-" gorm:"not null"` // "-" excludes from JSON
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
