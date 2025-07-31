package models

import (
	"time"

	"gorm.io/gorm"
)

// Folder represents a folder for organizing notes
type Folder struct {
	FolderID  uint           `json:"folderId" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"not null"`
	OwnerID   uint           `json:"ownerId" gorm:"not null;index"`
	Owner     User           `json:"owner" gorm:"foreignKey:OwnerID;references:UserID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// FolderShare represents folder sharing permissions
type FolderShare struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	FolderID  uint      `json:"folderId" gorm:"not null;index"`
	UserID    uint      `json:"userId" gorm:"not null;index"`
	Access    string    `json:"access" gorm:"not null;check:access IN ('read', 'write')"`
	Folder    Folder    `json:"folder" gorm:"foreignKey:FolderID;references:FolderID"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:UserID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName override for folder shares table
func (FolderShare) TableName() string {
	return "folder_shares"
}
