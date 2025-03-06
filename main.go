package main

import (
	"gin-jwt/controller"
	"gin-jwt/middleware"
	"gin-jwt/model"
	"os"
	"path"

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

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic("cache dir not found")
	}
	cacheDir = path.Join(cacheDir, "senaNoMusic")
	_, err = os.Stat(cacheDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(cacheDir, 0755)
		if err != nil {
			panic("error create caeh dir")
		}
	}
}

func main() {
	r := gin.Default()
	public := r.Group("/api")
	{
		// public.POST("/register", controller.Register)
		public.POST("/login", controller.Login)
	}

	apiV1 := r.Group("/api/v1")
	{
		apiV1.Use(middleware.JwtAuthMiddleware()) // 在路由组中使用中间件
		apiV1.GET("/user", controller.CurrentUser)
		apiV1.GET("/user/:id", controller.GetUserByID)
		apiV1.GET("/music/:id", controller.GetMusicById)
		apiV1.GET("/music/file/:id", controller.GetMusicStream)
		apiV1.GET("/music/search/:keyword", controller.ListMusicByTitle)
		apiV1.GET("/music/scan", controller.MusicScan)
		apiV1.GET("/music/list", controller.ListAllMusic)
		apiV1.GET("/music/clear", controller.ClearOldRecord)
	}

	apiV2 := r.Group("/api/v2")
	{
		apiV2.Use(middleware.JwtAuthMiddleware())
		apiV2.GET("/music/file/:id", controller.GetMusicStreamTrans)
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
