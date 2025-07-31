package teams

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTeam(db *gorm.DB, team *Team) (string, error) {
	team.ID = uuid.New().String()
	err := db.Create(team).Error
	return team.ID, err
}
