package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

const DB_PATH = "./Music.db"

// 连接数据库
func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open(DB_PATH), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&MusicInfo{})
}

func CreateAdminUser() (err error) {
	u := User{
		Username: "admin",
		Password: "test",
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
