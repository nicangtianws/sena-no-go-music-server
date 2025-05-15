package controller

import (
	"gin-jwt/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReqCollectInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MusicId     int    `json:"musicId"`
	CollectId   int    `json:"collectId"`
}

func CreateCollect(c *gin.Context) {
	var req ReqCollectInfo
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ResultErr(err.Error()))
		return
	}

	userId, err := GetUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultErr(err.Error()))
		return
	}

	model.CreateCollect(req.Name, req.Description, userId)

	c.JSON(http.StatusOK, ResultOk())

}

func AddMusicToCollect(c *gin.Context) {
	var req ReqCollectInfo
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ResultErr(err.Error()))
		return
	}

	userId, err := GetUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ResultErr(err.Error()))
		return
	}

	err = model.AddMusicToCollect(req.MusicId, req.CollectId, int(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultErr(err.Error()))
		return
	}

	c.JSON(http.StatusOK, ResultOk())
}

func DeleteMusicFromCollect(c *gin.Context) {
	var req ReqCollectInfo
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ResultErr(err.Error()))
		return
	}

	userId, err := GetUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ResultErr(err.Error()))
		return
	}

	err = model.DeleteMusicFromCollect(req.MusicId, req.CollectId, int(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResultErr(err.Error()))
		return
	}

	c.JSON(http.StatusOK, ResultOk())

}

func ListCollectByUserId(c *gin.Context) {
	userId, err := GetUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ResultErr(err.Error()))
		return
	}

	collectList := model.ListCollectByUserId(userId)

	c.JSON(http.StatusOK, ResultData(collectList))
}

func GetCollectByUserId(c *gin.Context) {
	idParam := c.Params.ByName("id")
	collectId, err := strconv.Atoi(idParam)
	if err != nil {
		c.String(http.StatusOK, "wrong id")
	}

	userId, err := GetUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ResultErr(err.Error()))
		return
	}

	collect := model.FindCollectById(collectId, userId)

	c.JSON(http.StatusOK, ResultData(collect))
}
