package controller

import (
	"gin-jwt/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req ReqUser

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": err.Error(),
		})
		return
	}

	u := model.User{
		Username: req.Username,
		Password: req.Password,
	}

	_, err := u.SaveUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "register success",
		"data":    req,
	})
}

func Login(c *gin.Context) {
	var req ReqUser
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := model.User{
		Username: req.Username,
		Password: req.Password,
	}

	// 对用户名和密码进行验证
	token, err := model.LoginCheck(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username or password is incorrect.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
