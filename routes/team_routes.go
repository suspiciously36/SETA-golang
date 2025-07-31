package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seta-namnv-6798/go-apis/controller"
)

// SetupTeamRoutes sets up all team-related routes
func SetupTeamRoutes(router *gin.Engine) {
	teamGroup := router.Group("/teams")
	{
		// Create a team
		teamGroup.POST("", controller.CreateTeam)

		// Team member management
		teamGroup.POST("/:teamId/members", controller.AddMemberToTeam)
		teamGroup.DELETE("/:teamId/members/:memberId", controller.RemoveMemberFromTeam)

		// Team manager management
		teamGroup.POST("/:teamId/managers", controller.AddManagerToTeam)
		teamGroup.DELETE("/:teamId/managers/:managerId", controller.RemoveManagerFromTeam)
	}
}
