package controller

import (
	"MFile/logic"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Likes(c *gin.Context) {
	type information struct {
		TargetName string `json:"targetName" binding:"required"`
		UserName   string `json:"userName" binding:"required"`
	}
	info := &information{}
	if err := c.ShouldBindJSON(info); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	err := logic.Likes(info.TargetName, info.UserName)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, "success")
	return
}
