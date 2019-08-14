package controller

import (
	"ahpuoj/request"
	"ahpuoj/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {

	var config = make(map[string]interface{})
	var enableIssueString string
	err := DB.Get(&enableIssueString, "select value from config where item = 'enable_issue'")
	if utils.CheckError(c, err, "数据库配置错误") != nil {
		return
	}

	if enableIssueString == "true" {
		config["enable_issue"] = true
	} else {
		config["enable_issue"] = false
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取系统配置项成功",
		"config":  config,
	})
}

func SetSettings(c *gin.Context) {
	var req request.Settings
	var config = make(map[string]interface{})
	c.ShouldBindJSON(&req)
	var enableIssueString string
	if req.EnableIssue == true {
		enableIssueString = "true"
	} else {
		enableIssueString = "false"
	}

	_, err := DB.Exec("update config set value = ? where item = 'enable_issue'", enableIssueString)
	if utils.CheckError(c, err, "数据库操作错误") != nil {
		return
	}
	if enableIssueString == "true" {
		config["enable_issue"] = true
	} else {
		config["enable_issue"] = false
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "修改系统配置项成功",
		"config":  config,
	})
}
