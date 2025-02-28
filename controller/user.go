package controller

import (
	"gin-jwt/model"
	"gin-jwt/utils/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
