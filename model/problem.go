package model

import (
	"ahpuoj/utils"
	"database/sql"
	"errors"
	"os"
	"strconv"
)

type Problem struct {
	Id           int            `db:"id"`
	Title        string         `db:"title"`
	Description  sql.NullString `db:"description"`
	Level        int            `db:"level"`
	Input        sql.NullString `db:"input"`
	Output       sql.NullString `db:"output"`
	SampleInput  sql.NullString `db:"sample_input"`
	SampleOutput sql.NullString `db:"sample_output"`
	Spj          int            `db:"spj"`
	Hint         sql.NullString `db:"hint"`
	Defunct      int            `db:"defunct"`
	TimeLimit    int            `db:"time_limit"`
	MemoryLimit  int            `db:"memory_limit"`
	Accepted     int            `db:"accepted"`
	Submit       int            `db:"submit"`
	Solved       int            `db:"solved"`
	CreatedAt    string         `db:"created_at"`
	UpdatedAt    string         `db:"updated_at"`
	CreatorId    string         `db:"creator_id"`
	Tags         []map[string]interface{}
}

func (problem *Problem) Save() error {
	result, err := DB.Exec(`insert into problem
	(title,description,input,output,sample_input,sample_output,spj,hint,level,time_limit,memory_limit,defunct,created_at,updated_at) 
	values (?,?,?,?,?,?,?,?,?,?,?,?,NOW(),NOW())`, problem.Title, problem.Description, problem.Input, problem.Output,
		problem.SampleInput, problem.SampleOutput, 0, problem.Hint, problem.Level, problem.TimeLimit, problem.MemoryLimit, 1)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	problem.Id = utils.Int64to32(lastInsertId)
	// 创建数据文件夹
	cfg := utils.GetCfg()
	dataDir, _ := cfg.GetValue("project", "datadir")
	baseDir := dataDir + "/" + strconv.FormatInt(int64(problem.Id), 10)
	utils.Consolelog(baseDir)
	err = os.MkdirAll(baseDir, 0777)
	return err
}

func (problem *Problem) Update() error {
	result, err := DB.Exec(`update problem 
	set title = ?,description = ?,input = ?,output = ?,sample_input = ?,sample_output = ?,spj = ?,
	hint = ?,level=?,time_limit = ?,memory_limit = ?,updated_at = NOW() where id = ?`, problem.Title, problem.Description, problem.Input, problem.Output,
		problem.SampleInput, problem.SampleOutput, 0, problem.Hint, problem.Level, problem.TimeLimit, problem.MemoryLimit, problem.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (problem *Problem) Delete() error {
	result, err := DB.Exec(`delete from problem where id = ?`, problem.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (problem *Problem) ToggleStatus() error {
	result, err := DB.Exec(`update problem set defunct = not defunct,updated_at = NOW() where id = ?`, problem.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (problem *Problem) AddTags(reqTags []interface{}) {

	tags := []map[string]interface{}{}
	insertStmt, _ := DB.Preparex(`insert into problem_tag (problem_id,tag_id,created_at,updated_at) values (?,?,NOW(),NOW())`)
	for _, tag := range reqTags {

		switch t := tag.(type) {
		// 如果是标签id 直接插入 golang解析json会将数字转成float64
		case float64:
			utils.Consolelog("float64")
			insertStmt.Exec(problem.Id, int(t))
			tag := Tag{
				Id: int(t),
			}
			tags = append(tags, tag.Response())
			// 否则生成新的tag并且插入
		case string:
			utils.Consolelog("string")
			newTag := Tag{
				Name: t,
			}
			err := newTag.Save()
			if err == nil {
				insertStmt.Exec(problem.Id, newTag.Id)
				tags = append(tags, newTag.Response())
			}
		}
	}
	insertStmt.Close()
	problem.Tags = tags
}

func (problem *Problem) FetchTags() {
	tags := []map[string]interface{}{}
	rows, err := DB.Queryx(`select tag.* from problem_tag inner join tag 
	on problem_tag.tag_id = tag.id where problem_tag.problem_id = ? order by problem_tag.id`, problem.Id)
	if err != nil {
		utils.Consolelog(err)
		return
	}
	for rows.Next() {
		utils.Consolelog("next")
		var tag Tag
		err = rows.StructScan(&tag)
		utils.Consolelog(tag)
		if err != nil {
			utils.Consolelog(err)
		}
		tags = append(tags, tag.Response())
	}
	rows.Close()
	problem.Tags = tags
}

func (problem *Problem) RemoveTags() error {
	_, err := DB.Exec(`delete from problem_tag where problem_id = ?`, problem.Id)
	problem.Tags = []map[string]interface{}{}
	return err
}

func (problem *Problem) Response() map[string]interface{} {

	return map[string]interface{}{
		"id":            problem.Id,
		"title":         problem.Title,
		"description":   problem.Description.String,
		"input":         problem.Input.String,
		"output":        problem.Output.String,
		"sample_input":  problem.SampleInput.String,
		"sample_output": problem.SampleOutput.String,
		"spj":           problem.Spj,
		"hint":          problem.Hint.String,
		"defunct":       problem.Defunct,
		"time_limit":    problem.TimeLimit,
		"memory_limit":  problem.MemoryLimit,
		"accepted":      problem.Accepted,
		"submit":        problem.Submit,
		"solved":        problem.Solved,
		"tags":          problem.Tags,
		"level":         problem.Level,
	}
}

func (problem *Problem) ResponseToUser() map[string]interface{} {
	// 需要将图片地址转换为绝对地址
	if problem.Description.Valid {
		problem.Description.String = utils.ConvertTextImgUrl(problem.Description.String)
	}
	if problem.Input.Valid {
		problem.Input.String = utils.ConvertTextImgUrl(problem.Input.String)
	}
	if problem.Output.Valid {
		problem.Output.String = utils.ConvertTextImgUrl(problem.Output.String)
	}
	if problem.Hint.Valid {
		problem.Hint.String = utils.ConvertTextImgUrl(problem.Hint.String)
	}
	return map[string]interface{}{
		"id":            problem.Id,
		"title":         problem.Title,
		"description":   problem.Description.String,
		"input":         problem.Input.String,
		"output":        problem.Output.String,
		"sample_input":  problem.SampleInput.String,
		"sample_output": problem.SampleOutput.String,
		"spj":           problem.Spj,
		"hint":          problem.Hint.String,
		"defunct":       problem.Defunct,
		"time_limit":    problem.TimeLimit,
		"memory_limit":  problem.MemoryLimit,
		"accepted":      problem.Accepted,
		"submit":        problem.Submit,
		"solved":        problem.Solved,
		"tags":          problem.Tags,
		"level":         problem.Level,
	}
}
