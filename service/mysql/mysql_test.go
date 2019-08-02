package mysql

import (
	"ahpuoj/utils"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDB(t *testing.T) {
	db, nil := NewDB()
	sql_check := "select tag_id from problem_tag where problem_id = 135"
	var count string

	var tags []int
	rows, err := db.Queryx(sql_check)
	for rows.Next() {
		utils.Consolelog("next")
		var tagId int
		err = rows.Scan(&tagId)
		if err != nil {
			utils.Consolelog(err)
		}
		tags = append(tags, tagId)
	}
	t.Log(tags)
	if err != nil {
		t.Log(err.Error())
	}
	t.Log(sql_check)
	t.Log(count)
	err = db.Ping()
	if err != nil {
		t.Errorf(err.Error())
		return
	} else {
	}
	defer db.Close()
}
