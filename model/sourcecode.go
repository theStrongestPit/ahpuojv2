package model

import (
	"errors"
)

type SourceCode struct {
	SolutionId int    `db:"soulution_id"`
	Source     string `db:"source"`
}

func (sourceCode *SourceCode) Save() error {
	_, err := DB.Exec(`insert into source_code
	(solution_id,source) values (?,?)`, sourceCode.SolutionId, sourceCode.Source)
	if err != nil {
		return err
	}
	return err
}

func (sourceCode *SourceCode) Delete() error {
	result, err := DB.Exec(`delete from source_code where soulution_id = ?`, sourceCode.SolutionId)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("数据不存在")
	}
	return err
}

func (sourceCode *SourceCode) Response() map[string]interface{} {

	return map[string]interface{}{
		"solution_id": sourceCode.SolutionId,
		"source":      sourceCode.Source,
	}
}
