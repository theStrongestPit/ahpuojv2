package model

import (
	"ahpuoj/utils"
	"errors"
)

type Series struct {
	Id          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	TeamMode    int    `db:"team_mode"`
	Defunct     int    `db:"defunct"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
	IsDeleted   int    `db:"is_deleted"`
}

func (series *Series) Save() error {
	result, err := DB.Exec(`insert into series
	(name,description,team_mode,created_at,updated_at) 
	values (?,?,?,NOW(),NOW())`, series.Name, series.Description, series.TeamMode)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	series.Id = utils.Int64to32(lastInsertId)
	return err
}

func (series *Series) Update() error {
	result, err := DB.Exec(`update series set name = ?,description=?,team_mode=?,updated_at = NOW() where id = ?`, series.Name, series.Description, series.TeamMode, series.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (series *Series) Delete() error {
	// 软删除
	result, err := DB.Exec(`update series set  is_deleted = 1 where id = ?`, series.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (series *Series) ToggleStatus() error {
	result, err := DB.Exec(`update series set defunct = not defunct,updated_at = NOW() where id = ?`, series.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (series *Series) Response() map[string]interface{} {

	return map[string]interface{}{
		"id":          series.Id,
		"name":        series.Name,
		"defunct":     series.Defunct,
		"description": series.Description,
		"team_mode":   series.TeamMode,
	}
}
