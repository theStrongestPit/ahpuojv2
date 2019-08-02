package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexNew(c *gin.Context) {

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
		whereString += "where title like '%" + param + "%'"
	}
	whereString += " order by top desc, id desc"
	utils.Consolelog(whereString)
	rows, total, err := model.Paginate(page, perpage, "new", []string{"*"}, whereString)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "数据获取失败",
		})
		return
	}
	var news []map[string]interface{}
	for rows.Next() {
		var new model.New
		rows.StructScan(&new)
		news = append(news, new.Response())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    news,
	})
}

func GetNew(c *gin.Context) {
	var new model.New
	id, _ := strconv.Atoi(c.Param("id"))
	stmt, _ := DB.Preparex("select * from new where id = ?")
	row := stmt.QueryRowx(id)
	if err := row.StructScan(&new); err != nil {
		c.JSON(400, gin.H{
			"message": "新闻不存在",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"new":     new.Response(),
	})
}

func StoreNew(c *gin.Context) {
	var req request.New
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "请求参数错误",
		})
		return
	}
	new := model.New{
		Title:   req.Title,
		Content: sql.NullString{req.Content, true},
	}
	err = new.Save()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "新建新闻失败，该新闻已存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "新建新闻成功",
		"new":     new.Response(),
	})
}

func UpdateNew(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.New
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "请求参数错误",
		})
		return
	}
	new := model.New{
		Id:      id,
		Title:   req.Title,
		Content: sql.NullString{req.Content, true},
	}
	err = new.Update()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "编辑新闻失败，该新闻已存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "编辑新闻成功",
		"new":     new.Response(),
	})
}

func DeleteNew(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	new := model.New{
		Id: id,
	}
	err := new.Delete()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "删除新闻失败，该新闻不存在",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "删除新闻成功",
	})
}

func ToggleNewStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	new := model.New{
		Id: id,
	}

	err := new.ToggleStatus()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "删除新闻失败，该新闻不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改新闻状态成功",
	})
}

func ToggleNewTopStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	new := model.New{
		Id: id,
	}

	err := new.ToggleTopStatus()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "更改新闻置顶状态失败，该新闻不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改新闻置顶状态成功",
	})
}
