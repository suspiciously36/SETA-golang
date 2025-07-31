package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seta-namnv-6798/go-apis/controller"
)

// SetupFolderRoutes sets up all folder-related routes
func SetupFolderRoutes(router *gin.Engine) {
	folderGroup := router.Group("/folders")
	{
		// Folder CRUD operations
		folderGroup.POST("", controller.CreateFolder)
		folderGroup.GET("/:folderId", controller.GetFolder)
		folderGroup.PUT("/:folderId", controller.UpdateFolder)
		folderGroup.DELETE("/:folderId", controller.DeleteFolder)

		// Folder sharing
		folderGroup.POST("/:folderId/share", controller.ShareFolder)
		folderGroup.DELETE("/:folderId/share/:userId", controller.RevokeFolderShare)

		// Notes within folders
		folderGroup.POST("/:folderId/notes", controller.CreateNote)
	}
}
