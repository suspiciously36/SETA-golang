package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seta-namnv-6798/go-apis/controller"
)

// SetupAssetRoutes sets up all asset management routes (Manager-only APIs)
func SetupAssetRoutes(router *gin.Engine) {
	// Team asset management
	teamGroup := router.Group("/teams")
	{
		teamGroup.GET("/:teamId/assets", controller.GetTeamAssets)
	}

	// User asset management
	userGroup := router.Group("/users")
	{
		userGroup.GET("/:userId/assets", controller.GetUserAssets)
	}
}
