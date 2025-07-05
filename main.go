package main

import (
	"flag"
	"gin-jwt/controller"
	"gin-jwt/middleware"
	"gin-jwt/model"
	"gin-jwt/util/audiofileutil"
	"gin-jwt/util/mylog"
	"os"

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

	// 获取基路径
	basedir := os.Getenv("DEFAULT_MUSIC_PATH")
	basedir = audiofileutil.AbsBasedir(basedir)

	// 初始化日志
	mylog.InitLog(&basedir)

	// 连初始化数据库
	model.InitDatabase(&basedir)

	// 初始化数据目录
	model.InitDir(&basedir)

	// 初始化管理员用户
	model.InitAdminUser()
}

func main() {
	r := gin.Default()
	public := r.Group("/api")
	{
		// public.POST("/register", controller.Register)
		public.POST("/login", controller.Login)
	}

	apiV1NoAuth := r.Group("/api/v1")
	{
		apiV1NoAuth.GET("/music/:id", controller.GetMusicById)
		apiV1NoAuth.GET("/music/file/:id", controller.GetMusicStream)
		apiV1NoAuth.GET("/music/search/:keyword", controller.ListMusicByTitle)
		apiV1NoAuth.GET("/music/list", controller.ListAllMusic)

	}

	apiV1 := r.Group("/api/v1")
	{
		apiV1.Use(middleware.JwtAuthMiddleware()) // 在路由组中使用中间件
		apiV1.GET("/user", controller.CurrentUser)
		apiV1.GET("/user/:id", controller.GetUserByID)
		apiV1.GET("/music/scan", controller.MusicScan)
		apiV1.GET("/music/clear", controller.ClearOldRecord)
		apiV1.POST("/collect", controller.CreateCollect)
		apiV1.GET("/collect", controller.ListCollectByUserId)
		apiV1.GET("/collect/:id", controller.GetCollectByUserId)
		apiV1.POST("/collect/add", controller.AddMusicToCollect)
		apiV1.POST("/collect/delete", controller.DeleteMusicFromCollect)
	}

	apiV2NoAuth := r.Group("/api/v2")
	{
		apiV2NoAuth.GET("/music/file/:id", controller.GetMusicStreamTrans)
	}

	apiV2 := r.Group("/api/v2")
	{
		apiV2.Use(middleware.JwtAuthMiddleware())
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
