package main

import (
	"gin-jwt/controller"
	"gin-jwt/middleware"
	"gin-jwt/model"
	"os"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	model.ConnectDatabase()
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	err = model.CreateAdminUser()
	if err != nil {
		panic("Error init admin user")
	}
}

func main() {
	r := gin.Default()
	public := r.Group("/api")
	{
		// public.POST("/register", controller.Register)
		public.POST("/login", controller.Login)
	}

	protected := r.Group("/api/admin")
	{
		protected.Use(middleware.JwtAuthMiddleware()) // 在路由组中使用中间件
		protected.GET("/user", controller.CurrentUser)
		protected.GET("/user/:id", controller.GetUserByID)
		protected.GET("/music/:id", controller.GetMusicById)
		protected.GET("/music/file/:id", controller.GetMusicStream)
		protected.GET("/music/search/:keyword", controller.ListMusicByTitle)
		protected.GET("/music/scan", controller.MusicScan)
		protected.GET("/music/list", controller.ListAllMusic)
		protected.GET("/music/clear", controller.ClearOldRecord)
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if strutil.IsBlank(host) {
		host = "127.0.0.1"
	}

	if strutil.IsBlank(port) {
		port = "8000"
	}

	r.Run(host + ":" + port)
}
