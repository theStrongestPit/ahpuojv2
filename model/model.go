package model

import (
	"ahpuoj/service/mysql"
	"ahpuoj/utils"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func Paginate(page int, perpage int, tableName string, outputField []string, whereString string) (*sqlx.Rows, int, error) {
	var total int
	var rows *sqlx.Rows
	offset := (page - 1) * perpage
	offsetStr := strconv.Itoa(offset)
	limitStr := strconv.Itoa(perpage)
	outputFieldString := strings.Join(outputField, ",")
	sql := "select " + outputFieldString + " from " + tableName + " " + whereString + " limit " + limitStr + " offset " + offsetStr
	utils.Consolelog(sql)
	rows, err := DB.Queryx(sql)
	if err != nil {
		utils.Consolelog(err)
	}
	sql = "select count(1) from (select 1 from " + tableName + " " + whereString + ")pesudo"
	utils.Consolelog(sql)
	DB.Get(&total, sql)
	return rows, total, err
}

func init() {
	DB, _ = mysql.NewDB()
	err := DB.Ping()
	if err != nil {
		utils.Consolelog(err.Error())
	} else {
		utils.Consolelog("successful connect to db")
	}
}
