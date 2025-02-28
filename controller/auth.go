package controller

import (
	"gin-jwt/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// /api/register的请求体
type ReqRegister struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req ReqRegister

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

// api/login 的请求体
type ReqLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req ReqLogin
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := model.User{
		Username: req.Username,
		Password: req.Password,
	}

	// 调用 models.LoginCheck 对用户名和密码进行验证
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
