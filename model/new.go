package model

import (
	"ahpuoj/utils"
	"database/sql"
	"errors"
)

type New struct {
	Id        int            `db:"id"`
	Title     string         `db:"title"`
	Content   sql.NullString `db:"content"`
	Top       int            `db:"top"`
	Defunct   int            `db:"defunct"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
}

func (new *New) Save() error {
	result, err := DB.Exec(`insert into new
	(title,content,top,defunct,created_at,updated_at) 
	values (?,?,0,0,NOW(),NOW())`, new.Title, new.Content)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	new.Id = utils.Int64to32(lastInsertId)
	return err
}

func (new *New) Update() error {
	result, err := DB.Exec(`update new set title = ?,content=?,updated_at = NOW() where id = ?`, new.Title, new.Content.String, new.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (new *New) Delete() error {
	result, err := DB.Exec(`delete from new where id = ?`, new.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (new *New) ToggleStatus() error {
	result, err := DB.Exec(`update new set defunct = not defunct,updated_at = NOW() where id = ?`, new.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (new *New) ToggleTopStatus() error {
	var newtop int
	// 查询top值
	DB.Get(&new.Top, `select top from new where id = ?`, new.Id)
	if new.Top == 0 {
		var maxtop int
		DB.Get(&maxtop, `select max(top) from new`)
		newtop = maxtop + 1
	} else {
		utils.Consolelog(123)
		newtop = 0
	}
	result, err := DB.Exec(`update new set top = ?, updated_at = NOW() where id = ?`, newtop, new.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (new *New) Response() map[string]interface{} {
	if new.Content.Valid {
		new.Content.String = utils.ConvertTextImgUrl(new.Content.String)
	}
	return map[string]interface{}{
		"id":         new.Id,
		"title":      new.Title,
		"content":    new.Content.String,
		"top":        new.Top,
		"defunct":    new.Defunct,
		"updated_at": new.UpdatedAt,
	}
}
