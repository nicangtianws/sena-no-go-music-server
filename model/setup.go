package model

import (
	"gin-jwt/util/mylog"
	"os"
	"path"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

const DB_FILE = "Music.db"

// 初始化数据库
func InitDatabase(basedir *string) {
	dbpath := filepath.Join(*basedir, DB_FILE)
	var err error
	DB, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&MusicInfo{})
	DB.AutoMigrate(&CollectInfo{})
}

// 初始化管理员
func InitAdminUser() {
	u := User{
		Username: os.Getenv("ADMIN_NAME"),
		Password: os.Getenv("ADMIN_PWD"),
	}

	existsUser := User{}

	DB.Model(&User{}).Where("username = ?", "admin").First(&existsUser)

	if existsUser.ID != 0 {
		return
	}
	_, err := u.SaveUser()
	if err != nil {
		mylog.LOG.Error().Msg("Error create admin user: " + err.Error())
		panic(err)
	}
}

// 初始化数据目录
func InitDir(basedir *string) {
	// 初始化缓存目录
	cacheDir := path.Join(*basedir, "cache")
	mylog.LOG.Info().Msg("Cache dir: " + cacheDir)
	_, err := os.Stat(cacheDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(cacheDir, 0755)
		if err != nil {
			panic("Error create cache dir: " + cacheDir)
		}
	}

	// 初始化存储目录
	musicDir := path.Join(*basedir, "music")
	mylog.LOG.Info().Msg("Music dir: " + musicDir)
	_, err = os.Stat(musicDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(musicDir, 0755)
		if err != nil {
			panic("Error create music dir: " + musicDir)
		}
	}
}
