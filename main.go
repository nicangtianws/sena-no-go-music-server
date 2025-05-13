package main

import (
	"flag"
	"gin-jwt/controller"
	"gin-jwt/middleware"
	"gin-jwt/model"
	"gin-jwt/utils/mylog"
	"os"
	"path"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// 初始化配置文件
	envFile := flag.String("env", ".env", "Path to the .env file")
	flag.Parse()
	err := godotenv.Load(*envFile)
	if err != nil {
		panic("Error loading .env file")
	}

	// 连初始化数据库
	model.ConnectDatabase()

	// 初始化日志
	mylog.LogInit()

	// 初始化管理员用户
	err = model.CreateAdminUser()
	if err != nil {
		panic("Error init admin user")
	}

	// 初始化缓存目录
	basedir := os.Getenv("DEFAULT_MUSIC_PATH")
	mylog.LOG.Info().Msg("base dir: " + basedir)
	cacheDir := path.Join(basedir, "cache")
	mylog.LOG.Info().Msg("cache dir: " + cacheDir)
	_, err = os.Stat(cacheDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(cacheDir, 0755)
		if err != nil {
			panic("error create cache dir: " + cacheDir)
		}
	}

	// 初始化存储目录
	musicDir := path.Join(basedir, "music")
	_, err = os.Stat(musicDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(musicDir, 0755)
		if err != nil {
			panic("error create music dir: " + musicDir)
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
		host = "0.0.0.0"
	}

	if strutil.IsBlank(port) {
		port = "8000"
	}

	r.Run(host + ":" + port)
}
