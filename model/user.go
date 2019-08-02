package model

import (
	"ahpuoj/utils"
	"database/sql"
	"errors"
)

type User struct {
	Id            int            `db:"id"`
	Email         sql.NullString `db:"email"`
	Username      string         `db:"username"`
	Nick          string         `db:"nick"`
	Avatar        string         `db:"avatar"`
	Password      string         `db:"password"`
	PassSalt      string         `db:"passsalt"`
	Submit        int            `db:"submit"`
	Solved        int            `db:"solved"`
	Defunct       int            `db:"defunct"`
	CreatedAt     string         `db:"created_at"`
	UpdatedAt     string         `db:"updated_at"`
	IsCompeteUser int            `db:"is_compete_user"`
	RoleId        int            `db:"role_id"`
	Role          string
}

func (user *User) Save() error {
	defaultAvatar, _ := utils.GetCfg().GetValue("preset", "avatar")
	result, err := DB.Exec(`insert into user
	(email,username,password,passsalt,nick,avatar,submit,solved,defunct,is_compete_user,created_at,updated_at) 
	values (?,?,?,?,?,?,0,0,0,?,NOW(),NOW())`, user.Email, user.Username, user.Password, user.PassSalt, user.Nick, defaultAvatar, user.IsCompeteUser)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	user.Id = utils.Int64to32(lastInsertId)
	return err
}

func (user *User) ToggleStatus() error {
	result, err := DB.Exec(`update user set defunct = not defunct,updated_at = NOW() where id = ?`, user.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (user *User) Update() error {
	result, err := DB.Exec(`update user set username = ?,updated_at = NOW() where id = ?`, user.Username, user.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (user *User) ChangePass() error {
	result, err := DB.Exec(`update user set password = ?, passsalt = ?,updated_at = NOW() where id = ?`, user.Password, user.PassSalt, user.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (user *User) Delete() error {
	result, err := DB.Exec(`delete from user where id = ?`, user.Id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (user *User) Response() map[string]interface{} {

	return map[string]interface{}{
		"id":       user.Id,
		"username": user.Username,
		"role":     user.Role,
		"nick":     user.Nick,
		"avatar":   user.Avatar,
		"submit":   user.Submit,
		"solved":   user.Solved,
		"defunct":  user.Defunct,
	}
}
