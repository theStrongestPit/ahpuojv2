package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"crypto/sha1"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 用户获取用户信息
func GetUser(c *gin.Context) {

	user, _ := c.Get("user")
	utils.Consolelog(user)
	if user, ok := user.(model.User); ok {
		c.JSON(200, gin.H{
			"message": "用户信息获取成功",
			"user":    user.Response(),
		})
	}
}

// 账号设置中重设昵称的接口
func ResetNick(c *gin.Context) {
	var user model.User
	userInterface, _ := c.Get("user")
	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}
	var req request.UserNick
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}
	user.Nick = req.Nick
	_, err = DB.Exec("update user set nick = ? where id = ?", req.Nick, user.Id)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "该昵称已被使用",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "昵称修改成功",
		"user":    user.Response(),
	})
}

// 账号设置中重设密码的接口
func ResetPassword(c *gin.Context) {
	var user model.User
	userInterface, _ := c.Get("user")
	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}
	var req request.UserResetPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	h := sha1.New()
	h.Write([]byte(user.PassSalt))
	h.Write([]byte(req.OldPassword))
	hashedOldPassword := fmt.Sprintf("%x", h.Sum(nil))

	if hashedOldPassword != user.Password {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "原密码错误",
		})
		return
	}

	// 更新密码
	// 加盐处理 16位随机字符串
	salt := utils.GetRandomString(16)
	h.Reset()
	h.Write([]byte(salt))
	h.Write([]byte(req.Password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))
	if hashedPassword == user.Password {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "新密码不能和原密码相同",
		})
		return
	}

	_, err = DB.Exec("update user set password = ?, passsalt = ? where id = ?", hashedPassword, salt, user.Id)
	utils.Consolelog(err)
	c.JSON(200, gin.H{
		"message": "密码修改成功",
	})
}

// 用户提交测试运行的接口
func SubmitToTestRun(c *gin.Context) {
	var err error

	user, _ := GetUserInstance(c)
	var req request.TestRun
	err = c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "提交失败,表单参数错误") != nil {
		return
	}

	// 提交记录
	solution := model.Solution{
		ProblemId:  0,
		TeamId:     0,
		UserId:     user.Id,
		ContestId:  0,
		Num:        0,
		IP:         c.ClientIP(),
		Language:   req.Language,
		CodeLength: len(req.Source),
	}
	err = solution.Save()
	if utils.CheckError(c, err, "保存提交记录失败") != nil {
		return
	}
	sourceCode := model.SourceCode{
		SolutionId: solution.Id,
		Source:     req.Source,
	}

	err = sourceCode.Save()
	if utils.CheckError(c, err, "保存代码记录失败") != nil {
		return
	}

	// 保存用户输入
	customInput := model.CustomInput{
		SolutionId: solution.Id,
		InputText:  req.InputText,
	}
	err = customInput.Save()
	if utils.CheckError(c, err, "保存用户输入失败") != nil {
		return
	}

	// 更新提交状态为等待评判
	_, err = DB.Exec("update solution set result = 0 where solution_id = ?", solution.Id)

	// 等待评测机评判
	var result int
	for {
		DB.Get(&result, "select  result from solution where solution_id = ?", solution.Id)
		if result >= 10 && result <= 13 {
			break
		}
		time.Sleep(time.Second)
	}

	// 获取结果
	var runtimeinfo string
	var compileinfo string
	var customOutput string

	err = DB.Get(&runtimeinfo, "select error from runtimeinfo where solution_id = ?", solution.Id)
	if err == nil {
		customOutput = runtimeinfo
	}
	err = DB.Get(&compileinfo, "select error from compileinfo where solution_id = ?", solution.Id)
	if err == nil {
		customOutput = compileinfo
	}

	// 删除测试运行的记录
	// DB.Exec("delect from solution where solution_id = ?", solution.Id)
	// DB.Exec("delect from runtimeinfo where solution_id = ?", solution.Id)
	// DB.Exec("delect from compileinfo where solution_id = ?", solution.Id)

	c.JSON(200, gin.H{
		"message":       "测试运行成功",
		"custom_output": customOutput,
	})
}

// 用户提交评测的接口
func SubmitToJudge(c *gin.Context) {
	var err error
	var contest model.Contest

	user, _ := GetUserInstance(c)
	var req request.Solution
	err = c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "提交失败") != nil {
		return
	}

	var problem model.Problem
	err = DB.Get(&problem, "select * from problem where id = ?", req.ProblemId)
	if utils.CheckError(c, err, "提交失败，问题不存在") != nil {
		return
	}

	submitable := false

	// 管理员无提交限制
	if user.Role == "admin" {
		submitable = true
	} else {
		// 比赛的提交
		if req.ContestId > 0 {

			err = DB.Get(&contest, "select * from contest where id = ? and is_deleted = 0", req.ContestId)
			if utils.CheckError(c, err, "提交失败，竞赛不存在") != nil {
				return
			}
			// 非管理员只有在比赛进行过程中并且有参加权限才能提交
			contest.CalcStatus()
			// 比赛进行中
			if contest.Status == 2 {
				// 公开
				if contest.Private == 0 {
					submitable = true
				} else {
					// 检测是否有提交权限
					var temp int
					DB.Get(&temp, "select count(1) from contest_user where contest_id = ? and user_id = ?", req.ContestId, user.Id)
					if temp > 0 {
						submitable = true
					}
				}
			}
		} else { // 平时的提交
			// 如果只是一般用户无法提交保留中的题目
			if problem.Defunct == 0 {
				submitable = true
			}
		}
	}

	if submitable {
		var teamId int
		// 如果为团队赛模式，并且非管理员提交，查询当前用户的teamId
		if contest.TeamMode == 1 && user.Role != "admin" {
			err = DB.Get(&teamId, "select team_id from contest_team_user ctu where ctu.contest_id = ? and ctu.user_id = ?", contest.Id, user.Id)
		}

		solution := model.Solution{
			ProblemId:  req.ProblemId,
			TeamId:     teamId,
			UserId:     user.Id,
			ContestId:  req.ContestId,
			Num:        req.Num,
			IP:         c.ClientIP(),
			Language:   req.Language,
			CodeLength: len(req.Source),
		}
		err := solution.Save()
		if utils.CheckError(c, err, "保存提交记录失败") != nil {
			return
		}
		sourceCode := model.SourceCode{
			SolutionId: solution.Id,
			Source:     req.Source,
		}
		err = sourceCode.Save()
		if utils.CheckError(c, err, "保存代码记录失败") != nil {
			return
		}

		// 更新提交状态为等待评判
		_, err = DB.Exec("update solution set result = 0 where solution_id = ?", solution.Id)
		c.JSON(200, gin.H{
			"message":  "提交成功",
			"solution": solution.Response(),
		})
	} else {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "对不起，你没有提交权限",
		})
	}
}

// 切换代码公开状态
func ToggleSolutionStatus(c *gin.Context) {
	var user model.User
	id, _ := strconv.Atoi(c.Param("id"))

	userInterface, _ := c.Get("user")
	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}

	var solutionUserId int
	DB.Get(&solutionUserId, "select user_id from solution where solution_id = ?", id)
	if user.Id != solutionUserId {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "对不起，你没有修改权限",
		})
		return
	}

	DB.Exec("update source_code set public = not public where solution_id = ?", id)
	c.JSON(200, gin.H{
		"message": "修改代码公开状态成功",
	})
}

// 下载题目数据文件
func DownloadDataFile(c *gin.Context) {
	var user model.User
	userInterface, _ := c.Get("user")
	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}

	pidStr := c.Query("pid")
	sidStr := c.Query("sid")
	filename := c.Query("filename")
	pid, _ := strconv.Atoi(pidStr)
	sid, _ := strconv.Atoi(sidStr)

	// 检验提交是否存在
	userId := 0
	err := DB.Get(&userId, "select 1 from solution where solution_id = ? and problem_id = ? and user_id = ?", sid, pid, user.Id)
	if err != nil {
		utils.Consolelog(err)
		c.AbortWithStatusJSON(400,
			gin.H{
				"message": "数据不存在",
			})
		return
	}
	// 检验错误信息是否与数据库信息匹配
	errFilename := ""
	err = DB.Get(&errFilename, "select error from runtimeinfo where solution_id = ?", sid)
	if err != nil {
		utils.Consolelog(err)
		c.AbortWithStatusJSON(400,
			gin.H{
				"message": "数据不存在",
			})
		return
	}
	errFilenameWithoutSuffix := strings.TrimSuffix(errFilename, filepath.Ext(errFilename))
	filenameWithoutSuffix := strings.TrimSuffix(filename, filepath.Ext(filename))
	if errFilenameWithoutSuffix != filenameWithoutSuffix {
		utils.Consolelog(err)
		c.AbortWithStatusJSON(400,
			gin.H{
				"message": "数据不存在",
			})
		return
	}

	// 读取文件
	cfg := utils.GetCfg()
	dataDir, _ := cfg.GetValue("project", "datadir")
	baseDir := dataDir + "/" + strconv.FormatInt(int64(pid), 10)
	dataFileName := baseDir + "/" + filename

	c.Header("Content-Disposition", `attachment; filename=`+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(dataFileName)
}

// 账号设置上传头像
func UploadAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("img")
	ext := path.Ext(header.Filename)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "头像上传失败",
		})
		return
	}
	url, err := utils.SaveFile(file, ext, "avatars")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "头像上传失败,请检查服务器设置",
		})
		return
	}
	var user model.User
	userInterface, _ := c.Get("user")
	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}
	// 如果不是默认头像 删除原头像
	defaultAvatar, _ := utils.GetCfg().GetValue("preset", "avatar")
	if user.Avatar != defaultAvatar {
		cfg := utils.GetCfg()
		webDir, _ := cfg.GetValue("project", "webdir")
		projectPath := webDir + "/"
		os.Remove(projectPath + user.Avatar)
	}
	DB.Exec("update user set avatar = ? where id = ?", url, user.Id)

	c.JSON(200, gin.H{
		"message": "头像上传成功",
		"url":     url,
	})
}

// 发布主题帖
func PostIssue(c *gin.Context) {
	var err error
	var user model.User

	userInterface, _ := c.Get("user")
	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}
	var req request.Issue
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	if req.ProblemId != 0 {
		var temp int
		DB.Get(&temp, "select count(1) from problem where id = ?", req.ProblemId)
		if temp == 0 {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "发布讨论主题失败，问题不存在",
			})
			return
		}
	}

	issue := model.Issue{
		Title:     req.Title,
		ProblemId: req.ProblemId,
		UserId:    user.Id,
	}

	err = issue.Save()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "发布讨论主题失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "发布讨论主题成功",
		"issue":   issue.Response(),
	})

}

// 回复主题帖
func ReplyToIssue(c *gin.Context) {
	var err error
	var user model.User
	issueId, _ := strconv.Atoi(c.Param("id"))

	userInterface, _ := c.Get("user")
	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}
	var req request.Reply
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	// 主题是否存在
	var temp int
	DB.Get(&temp, "select count(1) from issue where id = ?", issueId)
	if temp == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "发布回复失败，主题不存在",
		})
		return
	}

	// 如果是对回复的回复 检查该回复是否存在
	if req.ReplyId != 0 {
		DB.Get(&temp, "select count(1) from reply where id = ? and user_id = ?", req.ReplyId, req.ReplyUserId)
		if temp == 0 {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "发布回复失败，目标回复不存在",
			})
			return
		}
	}

	reply := model.Reply{
		IssueId:     issueId,
		UserId:      user.Id,
		ReplyId:     req.ReplyId,
		ReplyUserId: req.ReplyUserId,
		Content:     req.Content,
	}

	err = reply.Save()
	if err != nil {
		utils.Consolelog(err)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "发布回复失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "发布回复成功",
		"reply":   reply.Response(),
	})

}

// 获取最近一次提交的代码
func GetLatestSource(c *gin.Context) {
	var err error
	var user model.User

	userInterface, _ := c.Get("user")
	problemId, _ := strconv.Atoi(c.Param("id"))

	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}

	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	// 查找提交
	type SourceCode struct {
		Source   string `json:"source"`
		Language int    `json:"language"`
	}
	var sourceCode SourceCode
	err = DB.Get(&sourceCode, `select source_code.source,solution.language from solution inner join source_code 
	on source_code.solution_id = solution.solution_id where solution.problem_id = ? and solution.user_id = ? order by solution.in_date desc limit 1`, problemId, user.Id)
	c.JSON(200, gin.H{
		"message":    "获取最近提交信息成功",
		"sourcecode": sourceCode,
	})
}

// 获取最近一次比赛中提交的代码
func GetLatestContestSource(c *gin.Context) {
	var err error
	var user model.User

	userInterface, _ := c.Get("user")
	contestId, _ := strconv.Atoi(c.Param("id"))
	num, _ := strconv.Atoi(c.Param("num"))

	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}

	if err != nil {
		utils.Consolelog(err)
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	// 查找提交
	type SourceCode struct {
		Source   string `json:"source"`
		Language int    `json:"language"`
	}
	var sourceCode SourceCode
	err = DB.Get(&sourceCode, `select source_code.source,solution.language from solution inner join source_code 
	on source_code.solution_id = solution.solution_id where solution.contest_id = ? and solution.num = ? and solution.user_id = ? order by solution.in_date limit 1`, contestId, num, user.Id)
	c.JSON(200, gin.H{
		"message":    "获取最近提交信息成功",
		"sourcecode": sourceCode,
	})
}
