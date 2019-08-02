package model

import (
	"ahpuoj/utils"
	"database/sql"
)

type Solution struct {
	Id         int            `db:"solution_id"`
	ProblemId  int            `db:"problem_id"`
	TeamId     int            `db:"team_id"`
	UserId     int            `db:"user_id"`
	ContestId  int            `db:"contest_id"`
	Num        int            `db:"num"`
	Time       int            `db:"time"`
	Memory     int            `db:"memory"`
	InDate     string         `db:"in_date"`
	Result     int            `db:"result"`
	Language   int            `db:"language"`
	IP         string         `db:"ip"`
	JudgeTime  sql.NullString `db:"judgetime"`
	Valid      int            `db:"valid"`
	CodeLength int            `db:"code_length"`
	PassRate   float32        `db:"pass_rate"`
	LintError  int            `db:"lint_error"`
	Judger     sql.NullString `db:"judger"`
	// 附加信息
	Public       int    `db:"public"`
	Username     string `db:"username"`
	Nick         string `db:"nick"`
	UserAvatar   string `db:"avatar"`
	ProblemTitle string `db:"title"`
}

func (solution *Solution) Save() error {
	result, err := DB.Exec(`insert into solution
	(problem_id,team_id,user_id,contest_id,num,in_date,language,ip,code_length,result) 
	values (?,?,?,?,?,NOW(),?,?,?,14)`, solution.ProblemId, solution.TeamId, solution.UserId, solution.ContestId, solution.Num,
		solution.Language, solution.IP, solution.CodeLength)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	solution.Id = utils.Int64to32(lastInsertId)
	return err
}

func (solution *Solution) Response() map[string]interface{} {

	return map[string]interface{}{
		"id":          solution.Id,
		"team_id":     solution.TeamId,
		"contest_id":  solution.ContestId,
		"num":         solution.Num,
		"time":        solution.Time,
		"memory":      solution.Memory,
		"result":      solution.Result,
		"language":    solution.Language,
		"in_date":     solution.InDate,
		"judgetime":   solution.JudgeTime.String,
		"code_length": solution.CodeLength,
		"pass_rate":   solution.PassRate,
		"judger":      solution.Judger.String,
		"public":      solution.Public,
		"user": map[string]interface{}{
			"id":       solution.UserId,
			"username": solution.Username,
			"nick":     solution.Nick,
			"avatar":   solution.UserAvatar,
		},
		"problem": map[string]interface{}{
			"id":    solution.ProblemId,
			"title": solution.ProblemTitle,
		},
	}
}
