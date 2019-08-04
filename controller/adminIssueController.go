package controller

import (
	"ahpuoj/model"
	"ahpuoj/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ToggleIssueStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	issue := model.Issue{
		Id: id,
	}

	err := issue.ToggleStatus()
	if utils.CheckError(c, err, "更改主题状态失败，目标主题不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改用户状态成功",
	})
}

func ToggleReplyStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	reply := model.Reply{
		Id: id,
	}

	err := reply.ToggleStatus()
	if utils.CheckError(c, err, "更改回复状态失败，目标回复不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改用户状态成功",
	})
}
