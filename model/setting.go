package model

import "gorm.io/gorm"

type UserSetting struct {
	gorm.Model
	UserId uint `json:"userId"`
	Rate   uint `json:"rate"`
}

func SaveSetting() {

}
