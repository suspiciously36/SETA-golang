package teams

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	teamGroup := r.Group("/teams")
	{
		teamGroup.POST("", CreateTeamHandler(db))
		// bạn có thể thêm các API khác như thêm thành viên, thêm manager ở đây
	}
}
