package mysql

import (
	"ahpuoj/utils"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDB() (*sqlx.DB, error) {
	cfg := utils.GetCfg()
	dbcfg, _ := cfg.GetSection("mysql")
	path := strings.Join([]string{dbcfg["user"], ":", dbcfg["password"], "@tcp(", dbcfg["host"], ":", dbcfg["port"], ")/", dbcfg["database"], "?charset=utf8"}, "")
	db, _ := sqlx.Open("mysql", path)
	utils.Consolelog(db)
	err := db.Ping()
	return db, err
}
