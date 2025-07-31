package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seta-namnv-6798/go-apis/controller"
)

// SetupNoteRoutes sets up all note-related routes
func SetupNoteRoutes(router *gin.Engine) {
	noteGroup := router.Group("/notes")
	{
		// Note CRUD operations
		noteGroup.GET("/:noteId", controller.GetNote)
		noteGroup.PUT("/:noteId", controller.UpdateNote)
		noteGroup.DELETE("/:noteId", controller.DeleteNote)

		// Note sharing
		noteGroup.POST("/:noteId/share", controller.ShareNote)
		noteGroup.DELETE("/:noteId/share/:userId", controller.RevokeNoteShare)
	}
}
