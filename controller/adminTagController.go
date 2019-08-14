package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexTag(c *gin.Context) {

	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	param := c.Query("param")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	if page == 0 {
		page = 1
	}
	whereString := ""
	if len(param) > 0 {
		whereString += "where name like '%" + param + "%'"
	}
	whereString += " order by id desc"

	utils.Consolelog(whereString)
	rows, total, err := model.Paginate(page, perpage, "tag", []string{"*"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	var tags []map[string]interface{}
	for rows.Next() {
		var tag model.Tag
		rows.StructScan(&tag)
		tags = append(tags, tag.Response())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    tags,
	})
}

func GetAllTags(c *gin.Context) {
	rows, err := DB.Queryx("select * from tag order by id desc")
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	var tags []map[string]interface{}
	for rows.Next() {
		var tag model.Tag
		rows.StructScan(&tag)
		tags = append(tags, tag.Response())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"tags":    tags,
	})
}

func StoreTag(c *gin.Context) {
	var req request.Tag
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	tag := model.Tag{
		Name: req.Name,
	}
	err = tag.Save()
	if utils.CheckError(c, err, "新建标签失败，该标签已存在") != nil {
		return
	}

	c.JSON(200, gin.H{
		"message": "新建标签成功",
		"tag":     tag,
	})
}

func UpdateTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.Tag
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	tag := model.Tag{
		Id:   id,
		Name: req.Name,
	}
	err = tag.Update()
	if utils.CheckError(c, err, "编辑标签失败，该标签已存在") != nil {
		return
	}

	c.JSON(200, gin.H{
		"message": "编辑标签成功",
		"tag":     tag.Response(),
	})
}

func DeleteTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tag := model.Tag{
		Id: id,
	}
	err := tag.Delete()
	if utils.CheckError(c, err, "删除标签失败，该标签不存在") != nil {
		return
	}

	c.JSON(200, gin.H{
		"message": "删除标签成功",
	})
}
