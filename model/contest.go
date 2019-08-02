package model

import (
	"ahpuoj/utils"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Contest struct {
	Id           int            `db:"id"`
	Name         string         `db:"name"`
	StartTime    string         `db:"start_time"`
	EndTime      string         `db:"end_time"`
	Description  sql.NullString `db:"description"`
	Defunct      int            `db:"defunct"`
	Private      int            `db:"private"`
	TeamMode     int            `db:"team_mode"`
	LangMask     int            `db:"langmask"`
	CreatedAt    string         `db:"created_at"`
	UpdatedAt    string         `db:"updated_at"`
	IsDeleted    int            `db:"is_deleted"`
	CreatorId    string         `db:"creator_id"`
	Problems     string
	ProblemInfos []map[string]interface{}
	Status       int // 1代表未开始 2代表进行中 3代表已结束
}

func (contest *Contest) Save() error {
	result, err := DB.Exec(`insert into contest
	(name,description,start_time,end_time,defunct,private,team_mode,langmask,created_at,updated_at) 
	values (?,?,?,?,0,?,?,?,NOW(),NOW())`, contest.Name, contest.Description, contest.StartTime, contest.EndTime,
		contest.Private, contest.TeamMode, contest.LangMask)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	contest.Id = utils.Int64to32(lastInsertId)
	return err
}

func (contest *Contest) Update() error {

	var oldPrivate, oldTeamMode int
	DB.Get(&oldPrivate, "select private from contest where id = ?", contest.Id)
	DB.Get(&oldTeamMode, "select team_mode from contest where id = ?", contest.Id)

	result, err := DB.Exec(`update contest 
	set name = ?,description = ?,start_time = ?,end_time = ?,private = ?,team_mode = ?,langmask = ?,
	updated_at = NOW() where id = ?`, contest.Name, contest.Description, contest.StartTime, contest.EndTime,
		contest.Private, contest.TeamMode, contest.LangMask, contest.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}

	// 如果由私有模式更改为公开模式，或者团队模式变更，则重设人员列表
	if (oldPrivate == 1 && contest.Private == 0) || (oldTeamMode != contest.TeamMode) {
		utils.Consolelog("删")
		DB.Exec("delete from contest_user where contest_id = ?", contest.Id)
		DB.Exec("delete from contest_team where contest_id = ? ", contest.Id)
		DB.Exec("delete from contest_team_user where contest_id = ?", contest.Id)
	}

	return err
}

func (contest *Contest) Delete() error {
	// 软删除
	result, err := DB.Exec("update contest set is_deleted = 1 where id = ?", contest.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (contest *Contest) ToggleStatus() error {
	result, err := DB.Exec(`update contest set defunct = not defunct,updated_at = NOW() where id = ?`, contest.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (contest *Contest) AddProblems(reqProblems string) {

	pieces := strings.Split(reqProblems, ",")
	insertStmt, _ := DB.Preparex("insert into contest_problem (contest_id,problem_id, num, created_at,updated_at) values (?,?,?,NOW(),NOW())")
	checkStmt, _ := DB.Preparex("select id from problem where id = ?")
	if len(pieces) > 0 && len(pieces[0]) > 0 {
		for index, value := range pieces {
			problemId, _ := strconv.Atoi(value)
			utils.Consolelog(problemId)
			insertable := true
			var pid int
			err := checkStmt.Get(&pid, problemId)
			if err != nil {
				utils.Consolelog(err)
				insertable = false
			}
			if insertable {
				_, err = insertStmt.Exec(contest.Id, problemId, index+1)
				if err != nil {
					utils.Consolelog(err)
				}
			}
		}
	}
	insertStmt.Close()
	checkStmt.Close()
}

func (contest *Contest) FetchProblems() {
	problems := []string{}
	rows, err := DB.Queryx(`select problem_id from contest_problem where contest_id = ? order by num asc`, contest.Id)
	if err != nil {
		utils.Consolelog(err)
		return
	}
	for rows.Next() {
		var problemId int
		err = rows.Scan(&problemId)
		if err != nil {
			utils.Consolelog(err)
		}
		problems = append(problems, strconv.Itoa(problemId))
	}
	contest.Problems = strings.Join(problems, ",")
}

func (contest *Contest) AttachProblemInfo() {
	problemInfos := make([]map[string]interface{}, 0)
	rows, err := DB.Queryx(`select problem.* from contest_problem inner join problem on contest_problem.problem_id=problem.id 
	where contest_problem.contest_id = ? order by contest_problem.num asc`, contest.Id)
	if err != nil {
		utils.Consolelog(err)
		return
	}
	index := 1
	for rows.Next() {
		var problem Problem
		rows.StructScan(&problem)
		problemInfo := map[string]interface{}{
			"id":    problem.Id,
			"title": problem.Title,
			"num":   index,
		}
		index++
		problemInfos = append(problemInfos, problemInfo)
	}
	contest.ProblemInfos = problemInfos
}

func (contest *Contest) RemoveProblems() error {
	_, err := DB.Exec(`delete from contest_problem where contest_id = ?`, contest.Id)
	return err
}

// teamId大于0 则为向团队添加成员
func (contest *Contest) AddUsers(userlist string, teamId int) []string {
	pieces := strings.Split(userlist, "\n")
	var infos []string
	insertStmt, _ := DB.Preparex("insert into contest_user(contest_id,user_id,created_at,updated_at) VALUES (?,?,NOW(),NOW())")
	insertToTeamStmt, _ := DB.Preparex("insert into contest_team_user(contest_id,team_id,user_id,created_at,updated_at) VALUES (?,?,?,NOW(),NOW())")
	checkUserExistStmt, _ := DB.Preparex("select id from user where username = ?")
	checkHasUserStmt, _ := DB.Preparex("select 1 from contest_user where contest_user.contest_id = ? and contest_user.user_id = ?")
	checkUserBelongsToTeam, _ := DB.Preparex("select 1 from team_user where team_id = ? and user_id = ?")
	if len(pieces) > 0 && len(pieces[0]) > 0 {
		for _, username := range pieces {
			var userId, temp int
			var info string
			insertable := true
			// 判断用户是否存在 不存在则无法插入
			err := checkUserExistStmt.Get(&userId, username)
			if err != nil {
				insertable = false
				info = "竞赛&作业添加用户" + username + "失败，用户不存在"
			}
			// 判断是否已经添加了用户进入竞赛作业中
			err = checkHasUserStmt.Get(&temp, contest.Id, userId)
			// 有记录返回err==nil
			if userId > 0 && err == nil {
				insertable = false
				info = "竞赛&作业添加用户" + username + "失败，用户已被添加"
			}

			if teamId > 0 {
				err = checkUserBelongsToTeam.Get(&temp, teamId, userId)
				// 无记录返回err
				if userId > 0 && err != nil {
					insertable = false
					info = "竞赛&作业添加用户" + username + "失败，用户不属于该团队"
				}
			}
			utils.Consolelog(userId, contest.Id, teamId, username, insertable)
			if insertable {
				insertStmt.Exec(contest.Id, userId)

				// 向团队添加成员
				if teamId > 0 {
					insertToTeamStmt.Exec(contest.Id, teamId, userId)
				}
				info = "竞赛&作业添加用户" + username + "成功"
			}
			infos = append(infos, info)
		}
	}
	insertStmt.Close()
	insertToTeamStmt.Close()
	checkUserExistStmt.Close()
	checkHasUserStmt.Close()
	checkUserBelongsToTeam.Close()
	return infos
}

func (contest *Contest) CalcStatus() {
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", contest.StartTime, time.Local)
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", contest.EndTime, time.Local)
	if err != nil {
		return
	}
	nowTime := time.Now()
	if nowTime.Unix() < startTime.Unix() {
		contest.Status = 1
	} else if nowTime.Unix() > endTime.Unix() {
		contest.Status = 3
	} else {
		contest.Status = 2
	}
}

func (contest *Contest) Response() map[string]interface{} {
	// 需要将图片地址转换为绝对地址
	return map[string]interface{}{
		"id":          contest.Id,
		"name":        contest.Name,
		"start_time":  contest.StartTime,
		"end_time":    contest.EndTime,
		"problems":    contest.Problems,
		"description": contest.Description.String,
		"defunct":     contest.Defunct,
		"private":     contest.Private,
		"team_mode":   contest.TeamMode,
		"langmask":    contest.LangMask,
		"status":      contest.Status,
	}
}

func (contest *Contest) ResponseToUser() map[string]interface{} {
	// 需要将图片地址转换为绝对地址
	if contest.Description.Valid {
		contest.Description.String = utils.ConvertTextImgUrl(contest.Description.String)
	}
	return map[string]interface{}{
		"id":           contest.Id,
		"name":         contest.Name,
		"start_time":   contest.StartTime,
		"end_time":     contest.EndTime,
		"problems":     contest.Problems,
		"probleminfos": contest.ProblemInfos,
		"description":  contest.Description.String,
		"defunct":      contest.Defunct,
		"private":      contest.Private,
		"team_mode":    contest.TeamMode,
		"langmask":     contest.LangMask,
		"status":       contest.Status,
	}
}
func (contest *Contest) ListItemResponse() map[string]interface{} {

	return map[string]interface{}{
		"id":        contest.Id,
		"name":      contest.Name,
		"defunct":   contest.Defunct,
		"private":   contest.Private,
		"team_mode": contest.TeamMode,
		"status":    contest.Status,
	}
}
