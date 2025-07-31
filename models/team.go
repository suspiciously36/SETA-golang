package models

import (
	"time"

	"gorm.io/gorm"
)

// Team represents a team in the system
type Team struct {
	TeamID    uint           `json:"teamId" gorm:"primaryKey;autoIncrement"`
	TeamName  string         `json:"teamName" gorm:"not null"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TeamMember represents the many-to-many relationship between users and teams
type TeamMember struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"userId" gorm:"not null;index"`
	TeamID    uint      `json:"teamId" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:UserID"`
	Team      Team      `json:"team" gorm:"foreignKey:TeamID;references:TeamID"`
	CreatedAt time.Time `json:"createdAt"`
}

// TeamManager represents the many-to-many relationship between users and teams they manage
type TeamManager struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"userId" gorm:"not null;index"`
	TeamID    uint      `json:"teamId" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:UserID"`
	Team      Team      `json:"team" gorm:"foreignKey:TeamID;references:TeamID"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName overrides for mapping tables
func (TeamMember) TableName() string {
	return "team_members"
}

func (TeamManager) TableName() string {
	return "team_managers"
}
