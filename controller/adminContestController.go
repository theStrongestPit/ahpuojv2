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

func IndexContest(c *gin.Context) {

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

	rows, total, err := model.Paginate(page, perpage, "contest", []string{"*"}, whereString)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "数据获取失败",
		})
		return
	}
	var contests []map[string]interface{}
	for rows.Next() {
		var contest model.Contest
		rows.StructScan(&contest)
		utils.Consolelog(contest)
		contests = append(contests, contest.ListItemResponse())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    contests,
	})
}

func GetContest(c *gin.Context) {
	var contest model.Contest
	id, _ := strconv.Atoi(c.Param("id"))
	err := DB.Get(&contest, "select * from contest where id = ?", id)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "竞赛&作业不存在",
		})
		return
	}
	contest.FetchProblems()
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"contest": contest.Response(),
	})
}

func GetAllContests(c *gin.Context) {
	rows, _ := DB.Queryx("select * from contest where is_deleted = 0")
	var contests []map[string]interface{}
	for rows.Next() {
		var contest model.Contest
		rows.StructScan(&contest)
		contests = append(contests, contest.ListItemResponse())
	}
	c.JSON(200, gin.H{
		"message":  "数据获取成功",
		"contests": contests,
	})
}

func StoreContest(c *gin.Context) {
	var req request.Contest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "请求参数错误",
		})
		return
	}
	contest := model.Contest{
		Name:        req.Name,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Description: sql.NullString{req.Description, true},
		LangMask:    req.LangMask,
		Private:     req.Private,
		TeamMode:    req.TeamMode,
	}
	err = contest.Save()
	// 处理竞赛作业包含的问题
	contest.AddProblems(req.Problems)

	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "新建竞赛&作业失败，该竞赛&作业已存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "新建竞赛&作业成功",
		"contest": contest.Response(),
	})
}

func UpdateContest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.Contest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "请求参数错误",
		})
		return
	}
	contest := model.Contest{
		Id:          id,
		Name:        req.Name,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Description: sql.NullString{req.Description, true},
		LangMask:    req.LangMask,
		Private:     req.Private,
		TeamMode:    req.TeamMode,
	}

	err = contest.Update()
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "编辑竞赛&作业失败，竞赛&作业不存在",
		})
		return
	}

	// 处理题目列表
	contest.RemoveProblems()
	contest.AddProblems(req.Problems)

	c.JSON(200, gin.H{
		"message": "编辑竞赛&作业成功",
		"contest": contest.Response(),
	})
}

func DeleteContest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	contest := model.Contest{
		Id: id,
	}
	err := contest.Delete()
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "删除竞赛&作业失败，竞赛&作业不存在",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "删除竞赛&作业成功",
	})
}

func ToggleContestStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	contest := model.Contest{
		Id: id,
	}

	err := contest.ToggleStatus()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "更改竞赛&作业状态失败，竞赛&作业不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改竞赛&作业状态成功",
	})
}

// 处理个人赛人员列表
func IndexContestUser(c *gin.Context) {
	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	param := c.Query("param")
	contestId := c.Param("id")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	if page == 0 {
		page = 1
	}
	whereString := "where contest_user.contest_id=" + contestId
	if len(param) > 0 {
		whereString += " and user.username like '%" + param + "%' or user.nick like '%" + param + "%'"
	}
	whereString += " order by user.id desc"

	utils.Consolelog(whereString)
	rows, total, err := model.Paginate(page, perpage,
		"contest_user inner join user on contest_user.user_id = user.id",
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

func AddContestUsers(c *gin.Context) {
	var temp int
	var req request.ContestUsers
	id, _ := strconv.Atoi(c.Param("id"))
	c.ShouldBindJSON(&req)

	// 检查竞赛是否存在
	DB.Get(&temp, "select count(1) from contest where id = ? and is_deleted = 0", id)
	if temp == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "竞赛&作业不存在",
		})
		return
	}

	contest := model.Contest{
		Id: id,
	}
	infos := contest.AddUsers(req.UserList, 0)
	c.JSON(200, gin.H{
		"message": "操作成功",
		"info":    infos,
	})
}

func DeleteContestUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userId, _ := strconv.Atoi(c.Param("userid"))

	DB.Exec("delete from contest_user where contest_id = ? and user_id = ?", id, userId)
	c.JSON(200, gin.H{
		"message": "删除竞赛&作业人员成功",
	})
}

// 处理团队赛管理
func IndexContestTeamWithUser(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	rows, _ := DB.Queryx("select team.* from contest_team inner join team on contest_team.team_id = team.id where contest_team.contest_id = ?", id)
	teams := []map[string]interface{}{}
	for rows.Next() {
		var team model.Team
		rows.StructScan(&team)
		team.AttachUserInfo(id)
		teams = append(teams, team.ResponseWithUsers())
	}

	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"data":    teams,
	})
}

func AddContestTeam(c *gin.Context) {
	var err error
	var temp int
	id, _ := strconv.Atoi(c.Param("id"))
	teamId, _ := strconv.Atoi(c.Param("teamid"))
	// 检查竞赛是否存在
	DB.Get(&temp, "select count(1) from contest where id = ? and is_deleted = 0", id)
	if temp == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "竞赛&作业不存在",
		})
		return
	}

	// 检查团队是否存在
	DB.Get(&temp, "select count(1) from team where id = ? and is_deleted = 0", teamId)
	if temp == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "团队不存在",
		})
		return
	}

	// 检查是否已经添加进了竞赛作业中
	DB.Get(&temp, "select count(1) from contest_team where contest_id = ? and team_id = ? ", id, teamId)
	if temp > 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "该团队已经在该竞赛作业中",
		})
		return
	}

	_, err = DB.Exec("insert into contest_team(contest_id,team_id,created_at,updated_at) values(?,?,NOW(),NOW())", id, teamId)
	if err != nil {
		utils.Consolelog(err)
	}
	c.JSON(200, gin.H{
		"message": "添加团队成功",
	})
}

func DeleteContestTeam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	teamId, _ := strconv.Atoi(c.Param("teamid"))
	_, err := DB.Exec("delete from contest_team where contest_id = ? and team_id = ?", id, teamId)
	if err != nil {
		utils.Consolelog(err)
	}
	// 级联删除
	DB.Exec(`delete contest_user from contest_user inner join contest_team_user on contest_user.contest_id = contest_team_user.contest_id 
	where contest_user.contest_id = ? and contest_team_user.team_id = ?`, id, teamId)
	DB.Exec("delete from contest_team_user where contest_id = ? and team_id = ?", id, teamId)

	c.JSON(200, gin.H{
		"message": "删除团队成功",
	})
}

func DeleteContestTeamUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	teamId, _ := strconv.Atoi(c.Param("teamid"))
	userId, _ := strconv.Atoi(c.Param("userid"))
	DB.Exec(`delete contest_user from contest_user inner join contest_team_user on contest_user.contest_id = contest_team_user.contest_id 
	where contest_user.contest_id = ? and contest_user.user_id = ? and contest_team_user.team_id = ?`, id, userId, teamId)

	DB.Exec("delete from contest_team_user where contest_id = ? and team_id = ? and user_id = ?", id, teamId, userId)
	c.JSON(200, gin.H{
		"message": "删除团队人员成功",
	})
}

func AddContestTeamUsers(c *gin.Context) {
	var req request.ContestUsers
	var temp int
	id, _ := strconv.Atoi(c.Param("id"))
	teamId, _ := strconv.Atoi(c.Param("teamid"))
	c.ShouldBindJSON(&req)

	// 检查竞赛是否存在
	DB.Get(&temp, "select count(1) from contest where id = ? and is_deleted = 0", id)
	if temp == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "竞赛&作业不存在",
		})
		return
	}

	// 检查团队是否存在
	DB.Get(&temp, "select count(1) from team where id = ? and is_deleted = 0", teamId)
	if temp == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "团队不存在",
		})
		return
	}

	contest := model.Contest{
		Id: id,
	}

	infos := contest.AddUsers(req.UserList, teamId)

	c.JSON(200, gin.H{
		"message": "操作成功",
		"info":    infos,
	})
}

func AddContestTeamAllUsers(c *gin.Context) {
	var err error
	var temp int
	id, _ := strconv.Atoi(c.Param("id"))
	teamId, _ := strconv.Atoi(c.Param("teamid"))

	// 检查竞赛是否存在
	DB.Get(&temp, "select count(1) from contest where id = ? and is_deleted = 0", id)
	if temp == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "竞赛&作业不存在",
		})
		return
	}

	// 检查团队是否存在
	DB.Get(&temp, "select count(1) from team where id = ? and is_deleted = 0", teamId)
	if temp == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "团队不存在",
		})
		return
	}
	var infos []string
	rows, err := DB.Queryx("select user.* from user inner join team_user on user.id = team_user.user_id where team_user.team_id = ?", teamId)
	checkHasUserStmt, _ := DB.Preparex("select 1 from contest_user where contest_user.contest_id = ? and contest_user.user_id = ?")
	insertStmt, _ := DB.Preparex("insert into contest_user(contest_id,user_id,created_at,updated_at) VALUES (?,?,NOW(),NOW())")
	insertToTeamStmt, _ := DB.Preparex("insert into contest_team_user(contest_id,team_id,user_id,created_at,updated_at) VALUES (?,?,?,NOW(),NOW())")
	for rows.Next() {
		var user model.User
		var info string
		rows.StructScan(&user)
		utils.Consolelog(user)
		err = checkHasUserStmt.Get(&temp, id, user.Id)
		// 有记录返回err==nil
		if err == nil {
			info = "竞赛&作业添加用户" + user.Username + "失败，用户已被添加"
		} else {
			insertStmt.Exec(id, user.Id)
			insertToTeamStmt.Exec(id, teamId, user.Id)
			info = "竞赛&作业添加用户" + user.Username + "成功"
		}
		infos = append(infos, info)
	}
	insertStmt.Close()
	insertToTeamStmt.Close()
	checkHasUserStmt.Close()

	c.JSON(200, gin.H{
		"message": "操作成功",
		"info":    infos,
	})

}
