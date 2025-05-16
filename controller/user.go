package controller

import (
	"errors"
	"gin-jwt/model"
	"gin-jwt/util/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 请求体
type ReqUser struct {
	Id       int    `json:"Id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CurrentUser(c *gin.Context) {
	// 从token中解析出userId
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 根据userId从数据库查询数据
	u, err := model.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    u,
	})
}

// 从token中获取userId
func GetUserIdFromToken(c *gin.Context) (int, error) {
	// 从token中解析出userId
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		return -1, err
	}

	if userId < 1 {
		return -1, errors.New("invalid user")
	}

	return int(userId), err
}

func GetUserByID(c *gin.Context) {
	userId := c.Params.ByName("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.String(http.StatusOK, "wrong id")
	}

	u, err := model.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    u,
	})
}
