package model

import (
	"errors"
)

type CustomInput struct {
	SolutionId int    `db:"soulution_id"`
	InputText  string `db:"source"`
}

func (customInput *CustomInput) Save() error {
	_, err := DB.Exec(`insert into custominput
	(solution_id,input_text) values (?,?)`, customInput.SolutionId, customInput.InputText)
	if err != nil {
		return err
	}
	return err
}

func (customInput *CustomInput) Delete() error {
	result, err := DB.Exec(`delete from custominput where soulution_id = ?`, customInput.SolutionId)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (customInput *CustomInput) Response() map[string]interface{} {

	return map[string]interface{}{
		"solution_id": customInput.SolutionId,
		"input_text":  customInput.InputText,
	}
}
