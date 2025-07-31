package main

import (
	"log"

	"golang/internal/assets"
	"golang/internal/middleware"
	"golang/internal/shared"
	"golang/internal/teams"
	"golang/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	// Gọi hàm connect, không gán gì vì nó gán vào shared.DB
	shared.ConnectDatabase()

	r := gin.Default()

	// Middleware có thể dùng "*" khi chưa cần xác thực thực sự
	r.Use(middleware.AuthMiddleware("*"))

	teams.RegisterRoutes(r, shared.DB)
	assets.RegisterRoutes(r, shared.DB)
	user.RegisterRoutes(r, shared.DB)

	log.Fatal(r.Run(":8080"))
}
