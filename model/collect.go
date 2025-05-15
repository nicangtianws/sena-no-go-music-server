package model

import (
	"errors"
	"fmt"
	"gin-jwt/utils/mylog"

	"gorm.io/gorm"
)

// 音乐集
type CollectInfo struct {
	gorm.Model
	Id          int         `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	UserId      int         `json:"userId"`                                          // 所属用户
	Name        string      `json:"name"`                                            // 名称
	Description string      `json:"description"`                                     // 简介
	MusicInfos  []MusicInfo `json:"musicInfos" gorm:"many2many:collect_music_maps;"` // 关联音乐
}

type CollectMusicMap struct {
	MusicInfoID   uint `gorm:"column:music_info_id"`
	CollectInfoID uint `gorm:"column:collect_info_id"`
}

// @title 创建音乐集
// @param name string
// @param description string
func CreateCollect(name string, description string, userId int) {
	collect := CollectInfo{
		UserId:      userId,
		Name:        name,
		Description: description,
	}

	DB.Create(&collect)
}

func AddMusicToCollect(musicId int, collectId int, userId int) error {
	music := GetMusicById(&musicId)
	if music.Id < 1 {
		mylog.LOG.Error().Msg(fmt.Sprintf("valid music id: %d", musicId))
		return errors.New("music not found")
	}

	queryCollect := CollectInfo{
		UserId: userId,
		Id:     collectId,
	}

	DB.Find(&queryCollect).Preload("MusicInfos").First(&queryCollect)

	for _, musicInfo := range queryCollect.MusicInfos {
		if musicId == musicInfo.Id {
			mylog.LOG.Error().Msg(fmt.Sprintf("map of music and collect exists, musicId: %d, collectId: %d", musicId, collectId))
			return errors.New("map of music and collect exists")
		}
	}

	// if count > 0 {
	// 	mylog.LOG.Error().Msg(fmt.Sprintf("map of music and collect exists, musicId: %d, collectId: %d", musicId, collectId))
	// 	return errors.New("map of music and collect exists")
	// }

	collect := CollectInfo{
		Id:         collectId,
		UserId:     userId,
		MusicInfos: []MusicInfo{{Id: musicId}},
	}

	DB.Model(&collect).Updates(collect)

	return nil
}

func DeleteMusicFromCollect(musicId int, collectId int, userId int) error {
	queryCollect := CollectInfo{
		UserId: userId,
		Id:     collectId,
	}

	DB.Find(&queryCollect).Preload("MusicInfos").First(&queryCollect)

	isExist := false
	for _, musicInfo := range queryCollect.MusicInfos {
		if musicId == musicInfo.Id {
			isExist = true
		}
	}

	// 不在音乐集中直接返回成功
	if !isExist {
		return nil
	}

	DB.Delete(&CollectMusicMap{}, "music_info_id = ? AND collect_info_id = ?", musicId, collectId)

	return nil
}

func ListCollectByUserId(userId int) []CollectInfo {
	var collectList []CollectInfo
	DB.Find(&collectList, CollectInfo{UserId: userId})
	return collectList
}

func FindCollectById(collectId int, userId int) CollectInfo {
	queryCollect := CollectInfo{
		UserId: userId,
		Id:     collectId,
	}
	DB.Find(&queryCollect).Preload("MusicInfos").First(&queryCollect)

	return queryCollect
}
