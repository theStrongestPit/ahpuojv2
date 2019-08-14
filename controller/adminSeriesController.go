package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexSeries(c *gin.Context) {

	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	param := c.Query("param")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	if page == 0 {
		page = 1
	}
	whereString := " where is_deleted = 0 "
	if len(param) > 0 {
		whereString += "and name like '%" + param + "%'"
	}
	whereString += " order by id desc"

	utils.Consolelog(whereString)
	rows, total, err := model.Paginate(page, perpage, "series", []string{"*"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	var serieses []map[string]interface{}
	for rows.Next() {
		var series model.Series
		err = rows.StructScan(&series)
		utils.Consolelog(err)
		serieses = append(serieses, series.Response())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    serieses,
	})
}

func GetSeries(c *gin.Context) {
	var series model.Series
	id, _ := strconv.Atoi(c.Param("id"))
	err := DB.Get(&series, "select * from series where id = ?", id)
	if utils.CheckError(c, err, "系列赛不存在") != nil {
		return
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"series":  series.Response(),
	})
}

func IndexSeriesContest(c *gin.Context) {

	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	param := c.Query("param")
	seriesId := c.Param("id")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	if page == 0 {
		page = 1
	}
	whereString := "where contest_series.series_id=" + seriesId
	if len(param) > 0 {
		whereString += " and contest.name like '%" + param + "%'"
	}
	whereString += " order by contest.id desc"

	utils.Consolelog(whereString)
	rows, total, err := model.Paginate(page, perpage,
		"contest_series inner join contest on contest_series.contest_id = contest.id",
		[]string{"contest.*"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	contests := []map[string]interface{}{}
	for rows.Next() {
		var contest model.Contest
		rows.StructScan(&contest)
		utils.Consolelog(contest)
		contests = append(contests, contest.Response())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    contests,
	})
}

func StoreSeries(c *gin.Context) {
	var req request.Series
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	series := model.Series{
		Name:        req.Name,
		Description: req.Description,
		TeamMode:    req.TeamMode,
	}
	err = series.Save()
	if utils.CheckError(c, err, "新建系列赛失败，该系列赛已存在") != nil {
		return
	}

	c.JSON(200, gin.H{
		"message": "新建系列赛成功",
		"series":  series,
	})
}

func UpdateSeries(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.Series
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "请求参数错误") != nil {
		return
	}
	series := model.Series{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		TeamMode:    req.TeamMode,
	}
	err = series.Update()
	if utils.CheckError(c, err, "编辑系列赛失败，该系列赛已存在") != nil {
		return
	}

	c.JSON(200, gin.H{
		"message": "编辑系列赛成功",
		"series":  series.Response(),
	})
}

func ToggleSeriesStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	series := model.Series{
		Id: id,
	}

	err := series.ToggleStatus()
	if utils.CheckError(c, err, "更改竞赛&作业状态失败，竞赛&作业不存在") != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改竞赛&作业状态成功",
	})
}

func DeleteSeries(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	series := model.Series{
		Id: id,
	}
	err := series.Delete()
	if utils.CheckError(c, err, "删除系列赛失败，该系列赛不存在") != nil {
		return
	}

	c.JSON(200, gin.H{
		"message": "删除系列赛成功",
	})
}

func AddSeriesContest(c *gin.Context) {
	var err error
	var temp int
	id, _ := strconv.Atoi(c.Param("id"))
	contestId, _ := strconv.Atoi(c.Param("contestid"))
	// 检查系列赛是否存在
	DB.Get(&temp, "select count(1) from series where id = ? and is_deleted = 0", id)
	if temp == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "系列赛不存在",
		})
		return
	}

	// 检查竞赛&作业是否存在
	DB.Get(&temp, "select count(1) from contest where id = ? and is_deleted = 0", contestId)
	if temp == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "竞赛&作业不存在",
		})
		return
	}

	// 检查是否已经添加进了竞赛作业中
	DB.Get(&temp, "select count(1) from contest_series where series_id = ? and contest_id = ? ", id, contestId)
	if temp > 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "该竞赛&作业已经在系列赛中了",
		})
		return
	}
	_, err = DB.Exec("insert into contest_series(series_id,contest_id,created_at,updated_at) values(?,?,NOW(),NOW())", id, contestId)
	if utils.CheckError(c, err, "数据库操作失败") != nil {
		return
	}
	c.JSON(200, gin.H{
		"message": "添加竞赛&作业成功",
	})
}

func DeleteSeriesContest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	contestId, _ := strconv.Atoi(c.Param("contestid"))

	DB.Exec("delete from contest_series where series_id = ? and contest_id = ?", id, contestId)
	c.JSON(200, gin.H{
		"message": "删除系列赛竞赛&作业成功",
	})
}
