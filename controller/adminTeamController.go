package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTeam(c *gin.Context) {
	var team model.Team
	id, _ := strconv.Atoi(c.Param("id"))
	err := DB.Get(&team, "select * from team where id = ?", id)
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "团队不存在",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"team":    team.Response(),
	})
}

func IndexTeam(c *gin.Context) {

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
		whereString += "where name like '%" + param + "%'"
	}
	whereString += " order by id desc"

	utils.Consolelog(whereString)
	rows, total, err := model.Paginate(page, perpage, "team", []string{"*"}, whereString)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "数据获取失败",
		})
		return
	}
	var teams []map[string]interface{}
	for rows.Next() {
		var team model.Team
		rows.StructScan(&team)
		teams = append(teams, team.Response())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    teams,
	})
}

func GetAllTeams(c *gin.Context) {
	rows, _ := DB.Queryx("select * from team where is_deleted = 0 order by id desc")
	var teams []map[string]interface{}
	for rows.Next() {
		var team model.Team
		rows.StructScan(&team)
		teams = append(teams, team.Response())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"teams":   teams,
	})
}

func IndexTeamUser(c *gin.Context) {

	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	param := c.Query("param")
	teamId := c.Param("id")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	if page == 0 {
		page = 1
	}
	whereString := "where team_user.team_id=" + teamId
	if len(param) > 0 {
		whereString += " and user.username like '%" + param + "%' or user.nick like '%" + param + "%'"
	}
	whereString += " order by user.id desc"

	utils.Consolelog(whereString)
	rows, total, err := model.Paginate(page, perpage,
		"team_user inner join user on team_user.user_id = user.id",
		[]string{"user.*"}, whereString)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "数据获取失败",
		})
		return
	}
	users := []map[string]interface{}{}
	for rows.Next() {
		var user model.User
		rows.StructScan(&user)
		utils.Consolelog(user)
		users = append(users, user.Response())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    users,
	})
}

func AddTeamUsers(c *gin.Context) {
	var req request.TeamUsers
	id, _ := strconv.Atoi(c.Param("id"))
	c.ShouldBindJSON(&req)

	team := model.Team{
		Id: id,
	}
	utils.Consolelog(req.UserList)
	infos := team.AddUsers(req.UserList)

	c.JSON(200, gin.H{
		"message": "操作成功",
		"info":    infos,
	})
}

func StoreTeam(c *gin.Context) {
	var req request.Team
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "请求参数错误",
		})
		return
	}
	team := model.Team{
		Name: req.Name,
	}
	err = team.Save()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "新建团队失败，该团队已存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "新建团队成功",
		"team":    team,
	})
}

func UpdateTeam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.Team
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "请求参数错误",
		})
		return
	}
	team := model.Team{
		Id:   id,
		Name: req.Name,
	}
	err = team.Update()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "编辑团队失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "编辑团队成功",
		"team":    team.Response(),
	})
}

func DeleteTeam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	team := model.Team{
		Id: id,
	}
	err := team.Delete()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "删除团队失败，团队不存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "删除团队成功",
	})
}

func DeleteTeamUser(c *gin.Context) {
	teamId, _ := strconv.Atoi(c.Param("id"))
	userId, _ := strconv.Atoi(c.Param("userid"))

	result, _ := DB.Exec("delete from team_user where team_id = ? and user_id = ?", teamId, userId)

	// 级联删除
	_, err := DB.Exec(`delete contest_user from contest_user inner join contest_team_user on contest_user.contest_id = contest_team_user.contest_id 
	where contest_user.user_id = ? and contest_team_user.team_id = ?`, userId, teamId)
	if err != nil {
		utils.Consolelog(err)
	}
	_, err = DB.Exec("delete from contest_team_user where team_id = ? and user_id = ?", teamId, userId)
	if err != nil {
		utils.Consolelog(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(400, gin.H{
			"message": "删除团队成员失败，团队成员不存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "删除团队成员成功",
	})
}
