package controller

import (
	"gin-jwt/model"
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMusicStream(c *gin.Context) {
	musicId := c.Params.ByName("id")
	id, err := strconv.Atoi(musicId)
	if err != nil {
		c.String(http.StatusOK, "wrong id")
	}
	musicInfo := model.GetMusicById(&id)
	path := musicInfo.Path
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"music.mp3")
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(path)
}

func GetMusicStreamTrans(c *gin.Context) {
	musicId := c.Params.ByName("id")
	id, err := strconv.Atoi(musicId)
	if err != nil {
		c.String(http.StatusOK, "wrong id")
	}
	file := model.GetMusicTransFileById(&id)
	fileName := path.Base(file)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(file)
}

func GetMusicById(c *gin.Context) {
	musicId := c.Params.ByName("id")
	id, err := strconv.Atoi(musicId)
	if err != nil {
		c.String(http.StatusOK, "wrong id")
	}
	musicInfo := model.GetMusicById(&id)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    musicInfo,
	})
}

func ListMusicByTitle(c *gin.Context) {
	keyword := c.Params.ByName("keyword")
	musicList := model.FindMusicByTitle(&keyword)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    musicList,
	})
}

func MusicScan(c *gin.Context) {
	err := model.MusicScan()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultErr(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    "success",
	})
}

func ListAllMusic(c *gin.Context) {
	musicList := model.ListAllMusic()
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    musicList,
	})
}

func ClearOldRecord(c *gin.Context) {
	model.ClearOldRecord()
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    "success",
	})
}
