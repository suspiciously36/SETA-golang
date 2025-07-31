package user

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username string    `json:"username"`
	Email    string    `gorm:"unique" json:"email"`
	Role     string    `json:"role"` // "manager" or "member"
}
