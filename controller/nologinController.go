package controller

import (
	"ahpuoj/model"
	"ahpuoj/utils"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 访客获取新闻列表的接口
func NologinGetNewList(c *gin.Context) {

	var user model.User
	user, loggedIn := GetUserInstance(c)

	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	if page == 0 {
		page = 1
	}
	whereString := ""
	if loggedIn {
		if user.Role != "admin" {
			whereString += " where defunct = 0 "
		}
	}
	whereString += " order by top desc, id desc"

	rows, total, err := model.Paginate(page, perpage, "new", []string{"*"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
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
		"data":    news,
	})
}

// 访客获取问题列表的接口
func NologinGetProblemList(c *gin.Context) {

	var user model.User
	user, loggedIn := GetUserInstance(c)
	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	param := c.Query("param")
	levelStr := c.Query("level")
	tagIdStr := c.Query("tag_id")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	level, _ := strconv.Atoi(levelStr)
	tagId, _ := strconv.Atoi(tagIdStr)
	if page == 0 {
		page = 1
	}
	whereString := "where 1"
	if len(param) > 0 {
		_, err := strconv.Atoi(param)
		if err != nil {
			whereString += " and problem.title like '%" + param + "%'"
		} else {
			whereString += " and problem.id =" + param
		}
	}
	if tagId >= 0 {
		whereString += " and problem_tag.tag_id=" + tagIdStr
	}
	if level >= 0 {
		whereString += " and problem.level=" + levelStr
	}

	// 非管理员无法查看隐藏的题目
	if !loggedIn || (loggedIn && user.Role != "admin") {
		whereString += " and problem.defunct=0 "
	}

	whereString += " group by problem.id "
	whereString += " order by problem.id asc"
	utils.Consolelog(whereString)
	rows, total, err := model.Paginate(page, perpage, "problem left join problem_tag on problem.id = problem_tag.problem_id",
		[]string{"problem.*"}, whereString)
	utils.Consolelog(whereString)

	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}

	// 不统计比赛中的提交
	checkAcStmt, _ := DB.Preparex("select count(1) from solution where problem_id = ? and result = 4 and user_id = ? and contest_id = 0")
	defer checkAcStmt.Close()
	checkWaStmt, _ := DB.Preparex("select count(1) from solution where problem_id = ? and result != 4 and user_id = ? and contest_id = 0")
	defer checkWaStmt.Close()

	problems := make([]map[string]interface{}, 0)
	for rows.Next() {
		var problem model.Problem
		rows.StructScan(&problem)
		problem.FetchTags()

		status := 0
		var temp int
		// 如果登录了 查询ac信息
		if loggedIn {
			checkAcStmt.Get(&temp, problem.Id, user.Id)
			if temp > 0 {
				status = 1
			} else {
				checkWaStmt.Get(&temp, problem.Id, user.Id)
				if temp > 0 {
					status = 2
				}
			}
		}

		problems = append(problems, map[string]interface{}{
			"id":       problem.Id,
			"title":    problem.Title,
			"accepted": problem.Accepted,
			"submit":   problem.Submit,
			"solved":   problem.Solved,
			"tags":     problem.Tags,
			"level":    problem.Level,
			"status":   status,
		})
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    problems,
	})
}

// 访客获取竞赛列表的接口
func NologinGetContestList(c *gin.Context) {

	var user model.User

	user, loggedIn := GetUserInstance(c)

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
		whereString += " and name like '%" + param + "%' "
	}
	// 非管理员无法查看隐藏的竞赛
	if loggedIn {
		if user.Role != "admin" {
			whereString += " and defunct = 0 "
		}
	}

	whereString += " order by id desc"
	utils.Consolelog(whereString)
	rows, total, err := model.Paginate(page, perpage, "contest", []string{"*"}, whereString)

	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}

	var contests []map[string]interface{}
	for rows.Next() {
		var contest model.Contest
		rows.StructScan(&contest)
		contest.CalcStatus()
		contests = append(contests, contest.ResponseToUser())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    contests,
	})
}

// 访客获取评测记录列表的接口
func NologinGetSolutionList(c *gin.Context) {

	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	param := c.Query("param")
	username := c.Query("username")
	languageStr := c.Query("language")
	resultStr := c.Query("result")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	language, _ := strconv.Atoi(languageStr)
	result, _ := strconv.Atoi(resultStr)
	contestIdStr := c.Query("contest_id")
	contestId, _ := strconv.Atoi(contestIdStr)

	if page == 0 {
		page = 1
	}

	whereString := "where 1"
	if len(username) > 0 {
		whereString += " and (user.username='" + username + "' or user.nick='" + username + "')"
	}
	if language >= 0 {
		whereString += " and solution.language=" + languageStr
	}
	if result >= 0 {
		whereString += " and solution.result=" + resultStr
	}

	// 查询比赛中的提交
	if contestId > 0 {
		whereString += " and solution.contest_id=" + contestIdStr
		num, err := utils.EngNumToInt(param)
		// 搜索格式不对 直接PASS
		if utils.CheckError(c, err, "参数格式错误") != nil {
			return
		}

		if num > 0 {
			whereString += " and solution.num=" + strconv.Itoa(num)
		}
	} else {
		// 平时不显示比赛提交
		whereString += " and solution.contest_id=0 "
		if len(param) > 0 {
			_, err := strconv.Atoi(param)
			if err != nil {
				whereString += " and problem.title like '%" + param + "%'"
			} else {
				whereString += " and problem.id =" + param
			}
		}
	}

	whereString += " order by solution.solution_id desc"
	utils.Consolelog(whereString)
	// 多表联查
	rows, total, err := model.Paginate(page, perpage, `solution inner join problem on solution.problem_id=problem.id 
	inner join user on solution.user_id = user.id 
	inner join source_code on solution.solution_id=source_code.solution_id`,
		[]string{"solution.*,user.username,user.nick,user.avatar,problem.title,source_code.public"}, whereString)

	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}

	solutions := make([]map[string]interface{}, 0)
	for rows.Next() {
		var solution model.Solution
		rows.StructScan(&solution)
		solutions = append(solutions, solution.Response())
	}

	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    solutions,
	})
}

// 获取评测记录信息的接口
func NologinGetSolution(c *gin.Context) {
	var user model.User
	user, loggedIn := GetUserInstance(c)

	var solution model.Solution
	id, _ := strconv.Atoi(c.Param("id"))
	var err error
	err = DB.Get(&solution, `select solution.*,user.username,user.nick,user.avatar,problem.title,source_code.public from solution 
	inner join problem on solution.problem_id=problem.id 
	inner join user on solution.user_id = user.id 
	inner join source_code on solution.solution_id=source_code.solution_id 
	where solution.solution_id = ?`, id)

	if utils.CheckError(c, err, "记录不存在") != nil {
		return
	}

	seeable := false

	// 代码是否可以查看
	if loggedIn && user.Role == "admin" {
		seeable = true
	} else {
		// 自己的代码可以查看
		if loggedIn && solution.UserId == user.Id {
			seeable = true
		}
		// 非比赛中可以查看公开的代码
		if solution.ContestId == 0 {
			if solution.Public == 1 {
				seeable = true
			}
		}
	}

	var runtimeInfo string
	var compileInfo string
	var source string

	if seeable {
		DB.Get(&source, "select source from source_code where solution_id = ?", solution.Id)
	}
	// 当 result 为 WA TL ML PE  CE时，返回运行时错误信息
	if (solution.Result >= 5 && solution.Result <= 8) || solution.Result == 10 {
		DB.Get(&runtimeInfo, "select error from runtimeinfo where solution_id = ?", solution.Id)
	}
	if solution.Result == 11 {
		DB.Get(&compileInfo, "select error from compileinfo where solution_id = ?", solution.Id)
	}

	responseData := make(map[string]interface{}, 0)
	responseData["runtime_info"] = runtimeInfo
	responseData["compile_info"] = compileInfo
	responseData["source"] = source
	for k, v := range solution.Response() {
		responseData[k] = v
	}
	c.JSON(200, gin.H{
		"message":  "数据获取成功",
		"solution": responseData,
		"seeable":  seeable,
	})
}

// 获取全部标签的接口
func NologinGetAllTags(c *gin.Context) {
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

func NologinGetProblem(c *gin.Context) {
	var user model.User
	user, loggedIn := GetUserInstance(c)

	var problem model.Problem
	id, _ := strconv.Atoi(c.Param("id"))
	var err error

	if loggedIn && user.Role != "admin" || !loggedIn {
		err = DB.Get(&problem, "select * from problem where id = ? and defunct = 0", id)
	} else {
		err = DB.Get(&problem, "select * from problem where id = ?", id)
	}

	if utils.CheckError(c, err, "问题不存在") != nil {
		return
	}

	problem.FetchTags()

	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"problem": problem.ResponseToUser(),
	})
}

// 获取竞赛作业问题信息的接口
func NologinGetContestProblem(c *gin.Context) {
	var err error
	var user model.User
	user, loggedIn := GetUserInstance(c)
	cid, _ := strconv.Atoi(c.Param("id"))
	num, _ := strconv.Atoi(c.Param("num"))
	var contest model.Contest

	if loggedIn && user.Role != "admin" || !loggedIn {
		err = DB.Get(&contest, "select * from contest where id = ? and is_deleted = 0 and defunct = 0", cid)
	} else {
		err = DB.Get(&contest, "select * from contest where id = ? and is_deleted = 0", cid)
	}

	if utils.CheckError(c, err, "竞赛&作业不存在") != nil {
		return
	}

	var problem model.Problem
	contest.CalcStatus()
	seeable := true
	reason := ""

	if loggedIn {
		// 不是管理员
		if user.Role != "admin" {
			// 如果竞赛作业尚未开始，题目不可见
			if contest.Status == 1 {
				seeable = false
				reason = "竞赛尚未开始，题目不可见"
			} else if contest.Status == 3 { // 如果竞赛作业已经结束，题目可见
			} else { // 否则判断竞赛是否私有，私有判断是否有权限
				if contest.Private == 1 {
					var temp int
					DB.Get(&temp, "select count(1) from contest_user where contest_user.contest_id = ? and contest_user.user_id = ?", contest.Id, user.Id)
					if temp == 0 {
						seeable = false
						reason = "对不起，你没有参加此次竞赛的权限"
					}
				}
			}
		}
	} else { // 游客可以查看已经结束的竞赛的题目
		if contest.Status == 1 {
			seeable = false
			reason = "竞赛尚未开始，题目不可见"
		} else if contest.Status == 3 { // 如果竞赛作业已经结束，题目可见
		} else {
			if contest.Private == 1 { // 私有的竞赛作业无法查看
				seeable = false
				reason = "对不起，你没有参加此次竞赛的权限"
			}
		}
	}
	if seeable {
		err = DB.Get(&problem, `select problem.* from contest_problem inner join problem on contest_problem.problem_id = problem.id 
		where  contest_problem.contest_id= ? and contest_problem.num = ?`, cid, num)
		var contestSubmit, contestAccepted int
		// 处理提交和通过 只显示竞赛作业中的提交和通过 单人通过多次只算一次
		err = DB.Get(&contestSubmit, `select count(1) from solution where contest_id =  ? and num = ?`, cid, num)
		err = DB.Get(&contestAccepted, `select count(1)  from (select count(1) from solution where contest_id =  ? and num = ? and result = 4 group by user_id) T`, cid, num)
		problem.Submit = contestSubmit
		problem.Accepted = contestAccepted
	}

	if utils.CheckError(c, err, "问题不存在") != nil {
		return
	}

	if seeable {
		c.JSON(200, gin.H{
			"message": "数据获取成功",
			"seeable": seeable,
			"problem": problem.Response(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "数据获取成功",
			"seeable": seeable,
			"reason":  reason,
		})
	}
}

// 获取竞赛信息的接口
func NologinGetContest(c *gin.Context) {
	var err error
	var user model.User
	user, loggedIn := GetUserInstance(c)
	var contest model.Contest
	id, _ := strconv.Atoi(c.Param("id"))

	// 非管理员无法查看被保留的竞赛作业
	if loggedIn && user.Role != "admin" || !loggedIn {
		err = DB.Get(&contest, "select * from contest where id = ? and is_deleted = 0 and defunct = 0", id)
	} else {
		err = DB.Get(&contest, "select * from contest where id = ? and is_deleted = 0", id)
	}

	if utils.CheckError(c, err, "竞赛&作业不存在") != nil {
		return
	}

	contest.CalcStatus()
	seeable := true
	reason := ""

	if loggedIn {
		// 不是管理员
		if user.Role != "admin" {
			// 如果竞赛作业尚未开始，题目不可见
			if contest.Status == 1 {
				seeable = false
				reason = "竞赛尚未开始，题目不可见"
			} else if contest.Status == 3 { // 如果竞赛作业已经结束，题目可见
			} else { // 否则判断竞赛是否私有，私有判断是否有权限
				if contest.Private == 1 {
					var temp int
					DB.Get(&temp, "select count(1) from contest_user where contest_user.contest_id = ? and contest_user.user_id = ?", contest.Id, user.Id)
					if temp == 0 {
						seeable = false
						reason = "对不起，你没有参加此次竞赛的权限"
					}
				}
			}
		}
	} else { // 游客可以查看已经结束的竞赛的题目列表
		if contest.Status == 1 {
			seeable = false
			reason = "竞赛尚未开始，题目不可见"
		} else if contest.Status == 3 { // 如果竞赛作业已经结束，题目可见
		} else {
			if contest.Private == 1 { // 私有的竞赛作业无法查看
				seeable = false
				reason = "对不起，你没有参加此次竞赛的权限"
			}
		}
	}

	if seeable {
		contest.AttachProblemInfo()
		checkAcStmt, _ := DB.Preparex("select count(1) from solution where contest_id = ? and num = ? and result = 4 and user_id = ?")
		defer checkAcStmt.Close()
		checkWaStmt, _ := DB.Preparex("select count(1) from solution where contest_id = ? and num = ? and result != 4 and user_id = ?")
		defer checkWaStmt.Close()
		for k, v := range contest.ProblemInfos {
			contest.ProblemInfos[k]["status"] = 0
			// 如果已登录 检查AC状态 status = 1 通过 status = 2 错误 status = 0 未提交
			if loggedIn {
				problem := v
				var temp int
				checkAcStmt.Get(&temp, contest.Id, problem["num"], user.Id)
				if temp > 0 {
					contest.ProblemInfos[k]["status"] = 1
					continue
				}
				checkWaStmt.Get(&temp, contest.Id, problem["num"], user.Id)
				if temp > 0 {
					contest.ProblemInfos[k]["status"] = 2
					continue
				}
			}
		}

	} else {
		contest.ProblemInfos = make([]map[string]interface{}, 0)
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"seeable": seeable,
		"reason":  reason,
		"contest": contest.ResponseToUser(),
	})
}

// 获取竞赛作业排名的接口
func NologinGetContestRankList(c *gin.Context) {
	var user model.User
	user, loggedIn := GetUserInstance(c)
	var contest model.Contest
	id, _ := strconv.Atoi(c.Param("id"))

	err := DB.Get(&contest, "select * from contest where id = ? and is_deleted = 0", id)

	if utils.CheckError(c, err, "竞赛&作业不存在") != nil {
		return
	}

	contest.CalcStatus()
	seeable := true
	reason := ""

	if loggedIn {
		// 不是管理员
		if user.Role != "admin" {
			// 如果竞赛作业尚未开始，排名不可见
			if contest.Status == 1 {
				seeable = false
				reason = "竞赛尚未开始，排名不可见"
			} else if contest.Status == 3 { // 如果竞赛作业已经结束，排名可见
			} else { // 否则判断竞赛是否私有，私有判断是否有权限
				if contest.Private == 1 {
					var temp int
					DB.Get(&temp, "select count(1) from contest_user where contest_user.contest_id = ? and contest_user.user_id = ?", contest.Id, user.Id)
					if temp == 0 {
						seeable = false
						reason = "对不起，你没有参加此次竞赛的权限"
					}
				}
			}
		}
	} else { // 游客可以查看已经结束的竞赛的题目列表
		if contest.Status == 1 {
			seeable = false
			reason = "竞赛尚未开始，排名不可见"
		} else if contest.Status == 3 { // 如果竞赛作业已经结束，题目可见
		} else {
			if contest.Private == 1 { // 私有的竞赛作业无法查看
				seeable = false
				reason = "对不起，你没有参加此次竞赛的权限"
			}
		}
	}

	var userRankInfoList model.UserRankInfoList
	var problemCount int
	if seeable {
		// 获得竞赛作业题目总数
		DB.Get(&problemCount, "select count(1) from contest_problem where contest_id = ?", id)

		rows, _ := DB.Queryx(`select s.problem_id,s.team_id,s.user_id,s.contest_id,s.num,s.in_date,s.result,u.username,u.nick,u.avatar,r.name from
		solution s inner join user u on s.user_id = u.id 
		inner join role r on u.role_id = r.id
		where s.contest_id = ? order by s.user_id, s.in_date asc`, id)

		lastUserId := 0
		var userRankInfo model.UserRankInfo

		for rows.Next() {
			var rankItem model.RankItem
			rows.StructScan(&rankItem)
			// 忽略管理员的提交
			if rankItem.UserRole == "admin" {
				continue
			}

			// 如果是新的用户的数据
			if rankItem.UserId != lastUserId {
				if userRankInfo.User.Id != 0 {
					userRankInfoList = append(userRankInfoList, userRankInfo)
				}
				userRankInfo = model.UserRankInfo{
					Solved:  0,
					Time:    0,
					WaCount: make([]int, problemCount),
					AcTime:  make([]int, problemCount),
					User: struct {
						Id       int    `json:"id"`
						Username string `json:"username"`
						Nick     string `json:"nick"`
					}{
						Id:       rankItem.UserId,
						Username: rankItem.Username,
						Nick:     rankItem.Nick,
					},
				}
			}
			userRankInfo.Add(rankItem, contest.StartTime)
			lastUserId = rankItem.UserId
		}
		if userRankInfo.User.Id != 0 {
			userRankInfoList = append(userRankInfoList, userRankInfo)
		}
	}
	sort.Sort(userRankInfoList)
	c.JSON(200, gin.H{
		"message":  "数据获取成功",
		"seeable":  seeable,
		"reason":   reason,
		"ranklist": userRankInfoList,
		"contest": struct {
			ProblemCount int    `json:"problem_count"`
			Name         string `json:"name"`
			Id           int    `json:"id"`
		}{
			ProblemCount: problemCount,
			Name:         contest.Name,
			Id:           contest.Id,
		},
	})
}

// 获取竞赛作业团队排名的接口
func NologinGetContestTeamRankList(c *gin.Context) {
	var user model.User
	user, loggedIn := GetUserInstance(c)
	var contest model.Contest
	id, _ := strconv.Atoi(c.Param("id"))

	err := DB.Get(&contest, "select * from contest where id = ? and is_deleted = 0", id)

	if utils.CheckError(c, err, "竞赛&作业不存在") != nil {
		return
	}

	if contest.TeamMode != 1 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "竞赛&作业不是团队模式",
		})
		return
	}

	contest.CalcStatus()
	seeable := true
	reason := ""

	if loggedIn {
		// 不是管理员
		if user.Role != "admin" {
			// 如果竞赛作业尚未开始，排名不可见
			if contest.Status == 1 {
				seeable = false
				reason = "竞赛尚未开始，题目不可见"
			} else if contest.Status == 3 { // 如果竞赛作业已经结束，排名可见
			} else { // 否则判断竞赛是否私有，私有判断是否有权限
				if contest.Private == 1 {
					var temp int
					DB.Get(&temp, "select count(1) from contest_user where contest_user.contest_id = ? and contest_user.user_id = ?", contest.Id, user.Id)
					if temp == 0 {
						seeable = false
						reason = "对不起，你没有参加此次竞赛的权限"
					}
				}
			}
		}
	} else { // 游客可以查看已经结束的竞赛的题目列表
		if contest.Status == 1 {
			seeable = false
			reason = "竞赛尚未开始，排名不可见"
		} else if contest.Status == 3 { // 如果竞赛作业已经结束，题目可见
		} else {
			if contest.Private == 1 { // 私有的竞赛作业无法查看
				seeable = false
				reason = "对不起，你没有参加此次竞赛的权限"
			}
		}
	}

	// 按照team_id来排序
	var userRankInfoList model.UserRankSortByTeam
	var problemCount int
	if seeable {
		// 获得竞赛作业题目总数
		DB.Get(&problemCount, "select count(1) from contest_problem where contest_id = ?", id)
		rows, _ := DB.Queryx(`select s.problem_id,s.team_id,s.user_id,s.contest_id,s.num,s.in_date,s.result,u.username,u.nick,u.avatar,r.name from
		solution s inner join user u on s.user_id = u.id inner join role r on u.role_id = r.id where s.contest_id = ? order by s.user_id, s.in_date asc`, id)

		lastUserId := 0
		var userRankInfo model.UserRankInfo

		for rows.Next() {
			var rankItem model.RankItem
			rows.StructScan(&rankItem)
			// 忽略管理员的提交
			if rankItem.UserRole == "admin" {
				continue
			}
			// 如果是新的用户的数据
			if rankItem.UserId != lastUserId {
				if userRankInfo.User.Id != 0 {
					userRankInfoList = append(userRankInfoList, userRankInfo)
				}
				userRankInfo = model.UserRankInfo{
					Solved:  0,
					Time:    0,
					WaCount: make([]int, problemCount),
					AcTime:  make([]int, problemCount),
					TeamId:  rankItem.TeamId,
					User: struct {
						Id       int    `json:"id"`
						Username string `json:"username"`
						Nick     string `json:"nick"`
					}{
						Id:       rankItem.UserId,
						Username: rankItem.Username,
						Nick:     rankItem.Nick,
					},
				}
			}
			userRankInfo.TeamId = rankItem.TeamId
			userRankInfo.Add(rankItem, contest.StartTime)
			lastUserId = rankItem.UserId
		}
		userRankInfoList = append(userRankInfoList, userRankInfo)
	}
	sort.Sort(userRankInfoList)

	var teamRankInfoList model.TeamRankInfoList

	// 获取全部参赛队伍数据
	rows, _ := DB.Queryx(`select team.* from 
	contest_team inner join team on contest_team.team_id = team.id 
	where team.is_deleted = 0 and contest_team.contest_id = ? order by team.id asc`, contest.Id)
	for rows.Next() {
		var team model.Team
		rows.StructScan(&team)
		var teamRankInfo = model.TeamRankInfo{
			Solved:  0,
			Time:    0,
			WaCount: make([]int, problemCount),
			AcTime:  make([]int, problemCount),
			AcCount: make([]int, problemCount),
			Team: struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			}{
				Id:   team.Id,
				Name: team.Name,
			},
		}
		teamRankInfoList = append(teamRankInfoList, teamRankInfo)
	}

	// team排名信息和个人信息都是按照teamid递增排列的  o(n)方式来统计

	userCount := len(userRankInfoList)
	cnt := 0

out:
	for k, v := range teamRankInfoList {
		// 如果用户信息已经统计完 break
		if cnt >= userCount {
			break
		}
		utils.Consolelog(v.Team.Id, userRankInfoList[cnt].TeamId)
		// 如果个人信息的teamid大于当前team的id continue
		if userRankInfoList[cnt].TeamId > v.Team.Id {
			continue
		}
		for userRankInfoList[cnt].TeamId == v.Team.Id {
			teamRankInfoList[k].Add(userRankInfoList[cnt])
			cnt++
			if cnt >= userCount {
				break out
			}
		}
	}
	sort.Sort(teamRankInfoList)
	c.JSON(200, gin.H{
		"message":      "数据获取成功",
		"seeable":      seeable,
		"reason":       reason,
		"teamranklist": teamRankInfoList,
		"contest": struct {
			ProblemCount int    `json:"problem_count"`
			Name         string `json:"name"`
			Id           int    `json:"id"`
		}{
			ProblemCount: problemCount,
			Name:         contest.Name,
			Id:           contest.Id,
		},
	})
}

// 访客获取竞赛列表的接口
func NologinGetSeriesList(c *gin.Context) {

	var user model.User

	user, loggedIn := GetUserInstance(c)

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
		whereString += " and name like '%" + param + "%' "
	}
	// 非管理员无法查看隐藏的竞赛
	if loggedIn {
		if user.Role != "admin" {
			whereString += " and defunct = 0 "
		}
	}

	whereString += " order by id desc"
	rows, total, err := model.Paginate(page, perpage, "series", []string{"*"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}

	var serieses []map[string]interface{}
	for rows.Next() {
		var series model.Series
		rows.StructScan(&series)
		serieses = append(serieses, series.Response())
	}

	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    serieses,
	})
}

// 访客获取系列赛信息的接口  真难啊 写吐了都/(ㄒoㄒ)/~~
func NologinGetSeries(c *gin.Context) {
	var err error
	var user model.User
	user, loggedIn := GetUserInstance(c)
	var series model.Series
	id, _ := strconv.Atoi(c.Param("id"))

	if loggedIn && user.Role != "admin" || !loggedIn {
		err = DB.Get(&series, "select * from series where id = ? and defunct = 0", id)
	} else {
		err = DB.Get(&series, "select * from series where id = ?", id)
	}

	if utils.CheckError(c, err, "竞赛&作业不存在") != nil {
		return
	}

	series.AttachContestInfo()
	// 取得系列赛包含的竞赛作业数据
	rows, err := DB.Queryx("select contest.* from contest inner join contest_series on contest_series.contest_id = contest.id where contest_series.series_id = ? and contest.team_mode = ? and contest.is_deleted = 0 and contest.defunct = 0", id, series.TeamMode)
	if utils.CheckError(c, err, "数据库查询失败") != nil {
		return
	}
	var contestList []model.Contest
	contestStrList := ""
	for rows.Next() {
		var contest model.Contest
		rows.StructScan(&contest)
		contest.CalcStatus()
		contestList = append(contestList, contest)
		if len(contestStrList) > 0 {
			contestStrList += "," + strconv.Itoa(contest.Id)
		} else {
			contestStrList += strconv.Itoa(contest.Id)
		}
	}

	var contestCount int
	var problemCount int

	var userSeriesRankInfo model.UserSeriesRankInfo
	var userSeriesRankInfoList model.UserSeriesRankInfoList

	DB.Get(&contestCount, "select count(*) from contest_series inner join contest on contest_series.contest_id = contest.id where contest.is_deleted = 0 and contest.defunct = 0 and series_id = ? and contest.team_mode = ?", id, series.TeamMode)
	DB.Get(&problemCount, "select count(*) from contest_series inner join contest_problem on contest_series.contest_id = contest_problem.contest_id inner join contest on contest_series.contest_id = contest.id where contest.is_deleted = 0 and contest.defunct = 0 and series_id = ? and contest.team_mode = ?", id, series.TeamMode)

	// 如果竞赛作业数量为0，不进行后续处理
	if contestCount == 0 {
		c.JSON(200, gin.H{
			"message":      "数据获取成功",
			"series":       series.Response(),
			"userranklist": userSeriesRankInfoList,
		})
		return
	}

	// 个人模式排名汇总,取得系列赛全部的提交记录
	rows, err = DB.Queryx(`select s.problem_id,s.team_id,s.user_id,s.contest_id,s.num,s.in_date,s.result,u.username,u.nick,u.avatar,r.name from
	solution s inner join user u on s.user_id = u.id 
	inner join role r on u.role_id = r.id
	where s.contest_id in (` + contestStrList + `) order by s.user_id, s.in_date asc`)
	lastUserId := 0
	if utils.CheckError(c, err, "err") != nil {
		return
	}

	for rows.Next() {
		var rankItem model.RankItem
		var contest model.Contest

		rows.StructScan(&rankItem)
		// 获取当前提交的竞赛信息
		for _, c := range contestList {
			if rankItem.ContestId == c.Id {
				contest = c
				// break
			}
		}

		// 忽略管理员的提交
		if rankItem.UserRole == "admin" {
			continue
		}

		// 如果是新的用户的数据
		if rankItem.UserId != lastUserId {
			if userSeriesRankInfo.User.Id != 0 {
				userSeriesRankInfoList = append(userSeriesRankInfoList, userSeriesRankInfo)
			}
			userSeriesRankInfo = model.UserSeriesRankInfo{
				Solved:  make(map[int]int, contestCount),
				Time:    make(map[int]int, contestCount),
				WaCount: make(map[int][]int, problemCount),
				AcTime:  make(map[int][]int, problemCount),
				User: struct {
					Id       int    `json:"id"`
					Username string `json:"username"`
					Nick     string `json:"nick"`
				}{
					Id:       rankItem.UserId,
					Username: rankItem.Username,
					Nick:     rankItem.Nick,
				},
			}
		}
		userSeriesRankInfo.Add(rankItem, contest.Id, contest.StartTime, problemCount)
		lastUserId = rankItem.UserId
	}
	if userSeriesRankInfo.User.Id != 0 {
		userSeriesRankInfoList = append(userSeriesRankInfoList, userSeriesRankInfo)
	}

	// todolist 处理团队系列赛排名 这部分太复杂了 先搁置

	// 数据的排序交给前端处理，菜鸡不会用go处理这种排序(⊙﹏⊙)b

	c.JSON(200, gin.H{
		"message":      "数据获取成功",
		"series":       series.Response(),
		"userranklist": userSeriesRankInfoList,
	})
}

// 获取系统可用语言列表的接口
func NologinGetLanguageList(c *gin.Context) {
	cfg := utils.GetCfg()
	numberStr, _ := cfg.GetValue("language", "number")
	number, _ := strconv.Atoi(numberStr)
	langmaskStr, _ := cfg.GetValue("language", "mask")
	langmask, _ := strconv.Atoi(langmaskStr)
	langname, _ := cfg.GetValue("language", "langname")
	langNameList := strings.Split(langname, ",")
	languages := []map[string]interface{}{}
	for i := 0; i < number; i++ {
		available := false
		if (langmask & (1 << uint(i))) > 0 {
			available = true
		}
		lang := map[string]interface{}{
			"name":      langNameList[i],
			"available": available,
		}
		languages = append(languages, lang)
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "获取语言信息成功",
		"languages": languages,
	})
}

// 获取讨论列表的接口
func NologinGetIssueList(c *gin.Context) {

	// 判断当前是否已经关闭讨论版
	var enableIssueString string
	err := DB.Get(&enableIssueString, "select value from config where item = 'enable_issue'")
	if utils.CheckError(c, err, "数据库配置错误") != nil {
		return
	}

	if enableIssueString == "false" {
		c.JSON(200, gin.H{
			"message":      "数据获取成功",
			"issue_enable": false,
		})
		return
	}

	problemIdStr := c.Param("id")
	problemId, _ := strconv.Atoi(problemIdStr)

	user, loggedIn := GetUserInstance(c)
	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	if page == 0 {
		page = 1
	}

	// 检查问题是否存在
	if problemId != 0 {
		var temp int
		DB.Get(&temp, "select count(1) from problem where id = ?", problemId)
		if temp == 0 {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "问题不存在",
			})
			return
		}
	}

	whereString := "where 1"
	// problem=0时显示所有主题
	if problemId != 0 {
		whereString += " and problem_id =" + problemIdStr
	}
	// 管理员可以查看被删除的主题
	if !loggedIn || (loggedIn && user.Role != "admin") {
		whereString += " and is_deleted = 0 "
	}

	whereString += " order by created_at desc"

	rows, total, err := model.Paginate(page, perpage, "issue inner join user on issue.user_id = user.id left join problem on issue.problem_id = problem.id",
		[]string{"user.username,user.nick,user.avatar,issue.*,problem.title ptitle,(select count(1) from reply where issue_id = issue.id) as reply_count"}, whereString)
	utils.Consolelog(rows, total, err)

	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}

	var issues []map[string]interface{}
	for rows.Next() {
		var issue model.Issue
		err = rows.StructScan(&issue)
		utils.Consolelog(err)
		issues = append(issues, issue.Response())
	}
	c.JSON(200, gin.H{
		"message":      "数据获取成功",
		"total":        total,
		"issue_enable": true,
		"data":         issues,
	})
}

// 获得讨论以及评论的接口
func NologinGetIssue(c *gin.Context) {

	// 判断当前是否已经关闭讨论版
	var enableIssueString string
	err := DB.Get(&enableIssueString, "select value from config where item = 'enable_issue'")
	if utils.CheckError(c, err, "数据库配置错误") != nil {
		return
	}

	if enableIssueString == "false" {
		c.JSON(200, gin.H{
			"message":      "数据获取成功",
			"issue_enable": false,
		})
		return
	}

	issueIdStr := c.Param("id")
	issueId, _ := strconv.Atoi(issueIdStr)

	user, loggedIn := GetUserInstance(c)
	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	if page == 0 {
		page = 1
	}

	var issue model.Issue
	// 检查讨论是否存在
	if issueId != 0 {
		err = DB.Get(&issue, `select user.username,user.nick,user.avatar,issue.*,problem.title ptitle,(select count(1) from reply where issue_id = issue.id) as reply_count
		from issue inner join user on issue.user_id = user.id left join problem on issue.problem_id = problem.id where issue.id = ?`, issueId)
		if utils.CheckError(c, err, "讨论不存在") != nil {
			return
		}
	}

	// 第一步只获取对主题的回复
	whereString := "where issue_id = " + issueIdStr
	whereString += " and reply_id = 0"
	// 管理员可以查看被删除的回复
	if !loggedIn || (loggedIn && user.Role != "admin") {
		whereString += " and is_deleted = 0 "
	}

	whereString += " order by reply.created_at asc"

	rows, total, err := model.Paginate(page, perpage, "reply inner join user on reply.user_id = user.id",
		[]string{"user.username,user.nick,user.avatar,reply.*,'' as rnick,(select count(1) from reply  r where reply.id = r.reply_id) as reply_count"}, whereString)
	utils.Consolelog(rows, total, err)

	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}

	var replys []map[string]interface{}
	for rows.Next() {
		var reply model.Reply
		rows.StructScan(&reply)

		// 获取对该回复的回复
		if reply.ReplyCount > 0 {
			var subReplys []map[string]interface{}
			rows, _ := DB.Queryx(`select user.username,user.nick,user.avatar,reply.*,u2.nick as rnick,(select count(1) from reply  r where reply.id = r.reply_id) as reply_count
			from reply inner join user on reply.user_id = user.id inner join user u2 on reply.reply_user_id = u2.id where reply.reply_id = ?`, reply.Id)
			for rows.Next() {
				var subReply model.Reply
				err = rows.StructScan(&subReply)
				subReplys = append(subReplys, subReply.Response())
			}
			reply.SubReplys = subReplys
		}
		replys = append(replys, reply.Response())
	}
	c.JSON(200, gin.H{
		"message":      "数据获取成功",
		"total":        total,
		"issue":        issue.Response(),
		"replys":       replys,
		"issue_enable": true,
	})
}

// 获取用户信息的接口
func NologinGetUserInfo(c *gin.Context) {
	var err error
	userIdStr := c.Param("id")
	userId, _ := strconv.Atoi(userIdStr)

	var user model.User
	// 检查用户是否存在
	err = DB.Get(&user, `select * from user where id = ?`, userId)
	if utils.CheckError(c, err, "用户不存在") != nil {
		return
	}

	var solvedProblemList = make([]int, 0)
	var unsolvedProblemList = make([]int, 0)

	type StatisticUnit struct {
		Date  string `db:"date",json:"date"`
		Count int    `db:"count",json:"count"`
	}
	var recentSolvedStatistic = make([][]interface{}, 0)
	var recentSubmitStatistic = make([][]interface{}, 0)
	// 不统计比赛中的数据
	rows, _ := DB.Queryx("select distinct(problem_id) from solution where user_id=? and contest_id = 0 and result=4 and problem_id > 0 order by problem_id asc", userId)
	for rows.Next() {
		var pid int
		rows.Scan(&pid)
		solvedProblemList = append(solvedProblemList, pid)
	}

	rows, _ = DB.Queryx(`select distinct(problem_id) from solution where user_id=? and contest_id = 0 and result!=4 
	and problem_id not in (select distinct(problem_id) from solution where user_id=? and contest_id = 0 and result=4 order by problem_id asc)
		and problem_id > 0 order by problem_id asc`, user.Id, user.Id)

	for rows.Next() {
		var pid int
		rows.Scan(&pid)
		unsolvedProblemList = append(unsolvedProblemList, pid)
	}

	// 这是一段神奇的SQL 获得15天内累计通过的变化
	rows, _ = DB.Queryx(`
	select  dualdate.date,count(distinct(problem_id)) count from 
	(select * from solution where user_id=? and contest_id = 0 and result = 4) s 
	right join  
	(select date_sub(curdate(), interval(cast(help_topic_id as signed integer)) day) date
	from mysql.help_topic
	where help_topic_id  <= 14)  dualdate 
	on date(s.in_date) <= dualdate.date 
	group by dualdate.date order by dualdate.date asc`, user.Id)

	for rows.Next() {
		var unit StatisticUnit
		rows.StructScan(&unit)
		recentSolvedStatistic = append(recentSolvedStatistic, []interface{}{
			unit.Date,
			unit.Count,
		})
	}

	// 这还是一段神奇的SQL 获得15天内累计提交的变化
	rows, _ = DB.Queryx(`
	select  dualdate.date,count(distinct(problem_id)) count from 
	(select * from solution where user_id=? and contest_id = 0) s 
	right join  
	(select date_sub(curdate(), interval(cast(help_topic_id as signed integer)) day) date
	from mysql.help_topic
	where help_topic_id  <= 14)  dualdate 
	on date(s.in_date) <= dualdate.date 
	group by dualdate.date order by dualdate.date asc`, user.Id)

	for rows.Next() {
		var unit StatisticUnit
		rows.StructScan(&unit)
		recentSubmitStatistic = append(recentSubmitStatistic, []interface{}{
			unit.Date,
			unit.Count,
		})
	}
	var rank int

	DB.Get(&rank, "select count(1) from user where solved > ? or (solved = ? and submit < ?)", user.Solved, user.Solved, user.Submit)

	type UserInfo struct {
		Nick                  string          `json:"nick"`
		Avatar                string          `json:"avatar"`
		Solved                int             `json:"solved"`
		Submit                int             `json:"submit"`
		Rank                  int             `json:"rank"`
		CreatedAt             string          `json:"created_at"`
		SolvedProblemList     []int           `json:"solved_problem_list"`
		UnsolvedProblemList   []int           `json:"unsolved_problem_list"`
		RecentSolvedStatistic [][]interface{} `json:"recent_solved_statistic"`
		RecentSubmitStatistic [][]interface{} `json:"recent_submit_statistic"`
	}

	var userInfo UserInfo = UserInfo{
		Nick:                  user.Nick,
		Avatar:                user.Avatar,
		Solved:                user.Solved,
		Submit:                user.Submit,
		CreatedAt:             user.CreatedAt,
		Rank:                  rank + 1,
		SolvedProblemList:     solvedProblemList,
		UnsolvedProblemList:   unsolvedProblemList,
		RecentSolvedStatistic: recentSolvedStatistic,
		RecentSubmitStatistic: recentSubmitStatistic,
	}
	c.JSON(200, gin.H{
		"message":  "获取个人信息成功",
		"userinfo": userInfo,
	})
}

// 获取排名的接口
func NologinGetRankList(c *gin.Context) {
	pageStr := c.Query("page")
	perpageStr := "50"
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	if page == 0 {
		page = 1
	}

	// 不统计比赛用户的数据
	whereString := " where user.is_compete_user = 0 order by user.solved desc, user.submit asc"
	rows, total, err := model.Paginate(page, perpage, "user", []string{"user.*"}, whereString)
	if utils.CheckError(c, err, "数据获取失败") != nil {
		return
	}
	var users []map[string]interface{}
	for rows.Next() {
		var user model.User
		rows.StructScan(&user)
		users = append(users, user.Response())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    users,
	})
}
