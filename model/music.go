package model

import (
	"fmt"
	"gin-jwt/utils/audiofile"
	"gin-jwt/utils/ffmpegutil"
	"gin-jwt/utils/mylog"
	"os"
	"path"
	"path/filepath"

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
		mylog.LOG.Error().Msg(fmt.Sprintf("wrong path: %s", *path))
	}

	defer file.Close()

	title := file.Title()

	mylog.LOG.Info().Msg(fmt.Sprintf("title: %s", title))

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
		mylog.LOG.Error().Msg("cant find default path")
	}

	// homedir, _ := os.UserHomeDir()

	// dir := strings.Replace(basedir, "~", homedir, 1)

	dir := path.Join(basedir, "music")

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
			mylog.LOG.Warn().Msg("Not supported file type: " + err.Error())
		}

		if fileType.Is("audio/ogg") {
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
	mylog.LOG.Info().Msg(fmt.Sprintf("id: %d", *id))
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

func GetMusicTransFileById(id *int) string {
	mylog.LOG.Info().Msg(fmt.Sprintf("id: %d", *id))
	musicInfo := MusicInfo{Id: *id}
	DB.First(&musicInfo)
	sourcePath := musicInfo.Path

	fileType, err := mimetype.DetectFile(sourcePath)
	if err != nil {
		mylog.LOG.Warn().Msg("not supported file type: " + err.Error())
	}

	if !fileType.Is("audio/ogg") {
		basedir := os.Getenv("DEFAULT_MUSIC_PATH")
		cacheDir := path.Join(basedir, "cache")

		fileName := path.Base(sourcePath)

		cacheFile := cacheDir + "/" + fileName + "-320k.ogg"

		// 查看缓存文件是否存在，不存在则重新生成
		if !audiofile.CheckFileIsExist(cacheFile) {
			err = ffmpegutil.ConvertTo44kOGG(sourcePath, cacheFile)
			if err != nil {
				mylog.LOG.Error().Msg("convert failed: " + err.Error())
			}
		}

		return cacheFile
	}
	return sourcePath
}
