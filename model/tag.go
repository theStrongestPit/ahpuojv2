package model

import (
	"ahpuoj/utils"
	"errors"
)

type Tag struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (tag *Tag) Save() error {
	result, err := DB.Exec(`insert into tag
	(name,created_at,updated_at) 
	values (?,NOW(),NOW())`, tag.Name)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	tag.Id = utils.Int64to32(lastInsertId)
	return err
}

func (tag *Tag) Update() error {
	result, err := DB.Exec(`update tag set name = ?,updated_at = NOW() where id = ?`, tag.Name, tag.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (tag *Tag) Delete() error {
	// 级联删除
	_, err := DB.Exec(`delete from problem_tag where tag_id = ?`, tag.Id)
	if err != nil {
		return err
	}
	result, err := DB.Exec(`delete from tag where id = ?`, tag.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (tag *Tag) Response() map[string]interface{} {

	return map[string]interface{}{
		"id":   tag.Id,
		"name": tag.Name,
	}
}
