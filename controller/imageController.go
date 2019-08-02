package controller

import (
	"ahpuoj/utils"
	"path"

	"github.com/gin-gonic/gin"
)

func StoreImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	ext := path.Ext(header.Filename)
	if utils.CheckError(c, err, "图片获取失败") != nil {
		return
	}
	url, err := utils.SaveFile(file, ext, "images")
	if utils.CheckError(c, err, "图片保存失败") != nil {
		return
	}
	c.JSON(200, gin.H{
		"message": "图片上传成功",
		"url":     url,
	})
}
