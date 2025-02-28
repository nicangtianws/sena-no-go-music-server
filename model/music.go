package model

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gabriel-vasile/mimetype"
	"github.com/wtolson/go-taglib"
	"gorm.io/gorm"
)

type MusicInfo struct {
	gorm.Model
	Id         int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	BaseDir    string `json:"basedir"`
	Path       string `json:"path"`
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Album      string `json:"album"`
	Comment    string `json:"comment"`
	Genre      string `json:"genre"`
	Year       int    `json:"year"`
	Track      int    `json:"track"`
	Length     int    `json:"length"`
	Bitrate    int    `json:"bitrate"`
	Samplerate int    `json:"samplerate"`
	Channels   int    `json:"channels"`
}

func MusicParse(path *string, basedir *string) {
	file, err := taglib.Read(*path)

	if err != nil {
		slog.Info(fmt.Sprintf("wrong path: %s", *path))
	}

	defer file.Close()

	title := file.Title()

	slog.Info(fmt.Sprintf("title: %s", title))

	musicInfo := MusicInfo{
		Title:      title,
		Path:       *path,
		BaseDir:    *basedir,
		Artist:     file.Artist(),
		Album:      file.Album(),
		Comment:    file.Comment(),
		Genre:      file.Genre(),
		Year:       file.Year(),
		Track:      file.Track(),
		Length:     int(file.Length()),
		Bitrate:    file.Bitrate(),
		Samplerate: file.Samplerate(),
		Channels:   file.Channels(),
	}

	DB.Create(&musicInfo)
}

func MusicScan() {
	basedir := os.Getenv("DEFAULT_MUSIC_PATH")
	if strutil.IsBlank(basedir) {
		panic("cant find default path")
	}

	homedir, _ := os.UserHomeDir()

	dir := strings.Replace(basedir, "~", homedir, 1)

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 处理遍历过程中的错误
		}

		// 仅将文件路径加入到切片中
		if info.IsDir() {
			return nil
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return err // 获取绝对路径失败
		}

		fileType, err := mimetype.DetectFile(absPath)
		if err != nil {
			slog.Warn("Not supported file type: " + err.Error())
		}

		if fileType.Is("audio/ogg") || fileType.Is("audio/mpeg") {
			files = append(files, absPath)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	// 删除之前的记录重新扫描
	DB.Where("base_dir = ?", basedir).Delete(&MusicInfo{})
	for _, file := range files {
		// 根据path查找歌曲是否已经添加过
		musicList := FindMusicByPath(&file)
		if len(musicList) > 0 {
			continue
		}

		MusicParse(&file, &basedir)
	}
}

func FindMusicByPath(file *string) []MusicInfo {
	musicList := []MusicInfo{}
	DB.Where("path", file).Find(&musicList)
	return musicList
}

func GetMusicById(id *int) MusicInfo {
	slog.Info(fmt.Sprintf("id: %d", *id))
	musicInfo := MusicInfo{Id: *id}
	DB.First(&musicInfo)
	return musicInfo
}

func FindMusicByTitle(title *string) []MusicInfo {
	musicList := []MusicInfo{}
	titleLike := "%" + *title + "%"
	DB.Where("title LIKE ?", titleLike).Order("title asc").Find(&musicList)
	return musicList
}

func ListAllMusic() []MusicInfo {
	musicList := []MusicInfo{}
	DB.Model(&MusicInfo{}).Order("title asc").Find(&musicList)
	return musicList
}

func ClearOldRecord() {
	DB.Unscoped().Where("deleted_at is not null").Delete(&MusicInfo{})
}
