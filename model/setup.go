package model

import (
	"gin-jwt/utils/audiofileutil"
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

const DB_FILE = "Music.db"

// 连接数据库
func ConnectDatabase() {
	basedir := os.Getenv("DEFAULT_MUSIC_PATH")
	basedir = audiofileutil.AbsBasedir(basedir)

	dbpath := path.Join(basedir, DB_FILE)
	var err error
	DB, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&MusicInfo{})
	DB.AutoMigrate(&CollectInfo{})
}

func CreateAdminUser() (err error) {
	u := User{
		Username: os.Getenv("ADMIN_NAME"),
		Password: os.Getenv("ADMIN_PWD"),
	}

	existsUser := User{}

	DB.Model(&User{}).Where("username = ?", "admin").First(&existsUser)

	if existsUser.ID != 0 {
		return
	}
	_, err = u.SaveUser()
	if err != nil {
		panic(err)
	}
	return nil
}
