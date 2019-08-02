package controller

import (
	"ahpuoj/model"
	"ahpuoj/service/mysql"
	"ahpuoj/service/redisConn"
	"ahpuoj/utils"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB
var REDISPOOL *redis.Pool

func init() {
	DB, _ = mysql.NewDB()
	REDISPOOL = redisConn.NewPool()
	utils.Consolelog(REDISPOOL)

	err := DB.Ping()
	if err != nil {
		utils.Consolelog(err.Error())
	} else {
		utils.Consolelog("successful connect to db")
	}
}

// 获得user实例
func GetUserInstance(c *gin.Context) (model.User, bool) {
	var user model.User
	userInterface, loggedIn := c.Get("user")
	if userInterface, ok := userInterface.(model.User); ok {
		user = userInterface
	}
	return user, loggedIn
}
