package teams

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTeamHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input Team
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		input.ID = uuid.New().String()
		for i := range input.Managers {
			input.Managers[i].TeamID = input.ID
		}
		for i := range input.Members {
			input.Members[i].TeamID = input.ID
		}

		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
			return
		}

		c.JSON(http.StatusCreated, input)
	}
}
