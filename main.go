package main

import (
	"github.com/gin-gonic/gin"
	"github.com/seta-namnv-6798/go-apis/config"
	"github.com/seta-namnv-6798/go-apis/routes"
)

func main() {
	// Initialize database connection
	config.Connect()

	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Go APIs - Asset Management System")
	})

	// Setup all routes
	routes.SetupTeamRoutes(router)
	routes.SetupFolderRoutes(router)
	routes.SetupNoteRoutes(router)
	routes.SetupAssetRoutes(router)

	router.Run(":8080")
}
