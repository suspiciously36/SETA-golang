package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	users := r.Group("/users")
	{
		users.GET("/", GetAllUsers(db))
		users.POST("/", CreateUser(db))
	}
}
