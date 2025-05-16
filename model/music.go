package model

import (
	"fmt"
	"gin-jwt/util/audiofileutil"
	"gin-jwt/util/ffmpegutil"
	"gin-jwt/util/mylog"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/wtolson/go-taglib"
	"gorm.io/gorm"
)

// 基础信息
type MusicInfo struct {
	gorm.Model
	Id         int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	BaseDir    string `json:"basedir"`    // 所在文件夹
	Path       string `json:"path"`       // 绝对路径
	Title      string `json:"title"`      // 标题
	Artist     string `json:"artist"`     // 艺术家
	Album      string `json:"album"`      // 专辑
	Comment    string `json:"comment"`    // 简介
	Genre      string `json:"genre"`      // 风格
	Year       int    `json:"year"`       // 年份
	Track      int    `json:"track"`      // 轨道
	Length     int    `json:"length"`     // 时长
	Bitrate    int    `json:"bitrate"`    // 比特率
	Samplerate int    `json:"samplerate"` // 采样率
	Channels   int    `json:"channels"`   // 通道
}

func MusicParse(path *string, basedir *string) {
	file, err := taglib.Read(*path)

	if err != nil {
		mylog.LOG.Error().Msg(fmt.Sprintf("Wrong path: %s", *path))
		return
	}

	defer file.Close()

	title := file.Title()

	mylog.LOG.Info().Msg(fmt.Sprintf("Title: %s", title))

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

func MusicScan() error {
	basedir := os.Getenv("DEFAULT_MUSIC_PATH")
	basedir = audiofileutil.AbsBasedir(basedir)

	dir := path.Join(basedir, "music")

	// var musicList []MusicInfo

	var wg sync.WaitGroup
	fileChan := make(chan string, 100) // 文件处理队列

	// 创建worker
	numWorkers := runtime.NumCPU()
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				absPath, err := filepath.Abs(path)
				if err != nil {
					mylog.LOG.Warn().Msg("Cant get abs path: " + err.Error())
					return
				}

				// 是否在受支持音频类型列表内
				fileType, err := mimetype.DetectFile(absPath)
				if err != nil {
					mylog.LOG.Warn().Msg("Not supported file type: " + err.Error())
					return
				}

				_, err = audiofileutil.GetAudioFileType(fileType.String())

				if err != nil {
					mylog.LOG.Warn().Msg("Not supported file type: " + fileType.String())
					return
				}

				// files = append(files, absPath)
				// 根据path查找歌曲是否已经添加过
				musicList := FindMusicByPath(&path)
				if len(musicList) > 0 {
					continue
				}

				MusicParse(&path, &basedir)
			}
		}()
	}

	// 遍历目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 仅将文件路径加入到切片中
		if info.IsDir() {
			return nil
		}

		fileChan <- path

		return nil
	})
	if err != nil {
		return err
	}

	close(fileChan)
	wg.Wait()

	return nil

	// // 删除之前的记录重新扫描
	// DB.Where("base_dir = ?", basedir).Delete(&MusicInfo{})
	// for _, file := range files {
	// 	// 根据path查找歌曲是否已经添加过
	// 	musicList := FindMusicByPath(&file)
	// 	if len(musicList) > 0 {
	// 		continue
	// 	}

	// 	MusicParse(&file, &basedir)
	// }
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
	DB.Unscoped().Where("1=1").Delete(&MusicInfo{})
}

func GetMusicTransFileById(id *int) string {
	mylog.LOG.Info().Msg(fmt.Sprintf("music id: %d", *id))
	musicInfo := MusicInfo{Id: *id}
	DB.First(&musicInfo)
	sourcePath := musicInfo.Path

	fileType, err := mimetype.DetectFile(sourcePath)
	if err != nil {
		mylog.LOG.Warn().Msg("not supported file type: " + err.Error())
		return ""
	}

	// mp3不进行转换
	if fileType.Is("audio/mp3") {
		return sourcePath
	}

	// 其他类型如wav、flac需要进行转换
	if !fileType.Is("audio/ogg") {
		basedir := os.Getenv("DEFAULT_MUSIC_PATH")
		basedir = audiofileutil.AbsBasedir(basedir)
		cacheDir := path.Join(basedir, "cache")

		fileName := path.Base(sourcePath)

		cacheFile := path.Join(cacheDir, fileName+"-320k.ogg")
		coverFile := path.Join(cacheDir, fileName+"-320k.jpg")

		// 查看缓存文件是否存在，不存在则重新生成
		if !audiofileutil.CheckFileIsExist(cacheFile) {
			err = ffmpegutil.ConvertTo44kOGG(sourcePath, cacheFile)
			if err != nil {
				// 移除生成失败的文件
				os.Remove(cacheFile)
				mylog.LOG.Err(err)
			}
		}

		// 提取封面
		if !audiofileutil.CheckFileIsExist(coverFile) {
			err = ffmpegutil.ConvertCover(sourcePath, coverFile)
			if err != nil {
				// 移除生成失败的文件
				os.Remove(coverFile)
				mylog.LOG.Err(err)
			}
		}

		return cacheFile
	}
	return sourcePath
}
