package models

import (
	"time"

	"gorm.io/gorm"
)

// Note represents a note in the system
type Note struct {
	NoteID    uint           `json:"noteId" gorm:"primaryKey;autoIncrement"`
	Title     string         `json:"title" gorm:"not null"`
	Body      string         `json:"body" gorm:"type:text"`
	FolderID  uint           `json:"folderId" gorm:"not null;index"`
	OwnerID   uint           `json:"ownerId" gorm:"not null;index"`
	Folder    Folder         `json:"folder" gorm:"foreignKey:FolderID;references:FolderID"`
	Owner     User           `json:"owner" gorm:"foreignKey:OwnerID;references:UserID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// NoteShare represents note sharing permissions
type NoteShare struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	NoteID    uint      `json:"noteId" gorm:"not null;index"`
	UserID    uint      `json:"userId" gorm:"not null;index"`
	Access    string    `json:"access" gorm:"not null;check:access IN ('read', 'write')"`
	Note      Note      `json:"note" gorm:"foreignKey:NoteID;references:NoteID"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:UserID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName override for note shares table
func (NoteShare) TableName() string {
	return "note_shares"
}
