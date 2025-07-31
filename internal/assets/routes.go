package assets

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// tạo group /assets
	assetsGroup := r.Group("/assets")

	// thêm các route API (GET, POST, v.v.)
	assetsGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong from assets"})
	})

	// TODO: add real endpoints (folders, notes...) later
}
