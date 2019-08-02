package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StoreCasbin(c *gin.Context) {
	var req request.Casbin
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "请求参数错误",
		})
		return
	}
	ptype := "p"
	casbin := model.Casbin{
		Ptype:    ptype,
		RoleName: req.Rolename,
		Path:     req.Path,
		Method:   req.Method,
	}
	err = casbin.Store()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "新建权限规则失败",
		})

	}
	c.JSON(http.StatusOK, gin.H{
		"message": "新建权限规则成功",
	})
}
