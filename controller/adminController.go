package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetSubmitStatistic(c *gin.Context) {
	var rows *sqlx.Rows

	// 这还是一段神奇的SQL 获得15天内累计提交的变化
	type StatisticUnit struct {
		Date  string `db:"date",json:"date"`
		Count int    `db:"count",json:"count"`
	}
	var recentSubmitStatistic = make([][]interface{}, 0)

	rows, _ = DB.Queryx(`
	select  dualdate.date,count(*) count from 
	(select * from solution) s 
	right join  
	(select date_sub(curdate(), interval(cast(help_topic_id as signed integer)) day) date
	from mysql.help_topic
	where help_topic_id  <= 14)  dualdate 
	on date(s.in_date) <= dualdate.date 
	group by dualdate.date order by dualdate.date asc`)

	for rows.Next() {
		var unit StatisticUnit
		rows.StructScan(&unit)
		recentSubmitStatistic = append(recentSubmitStatistic, []interface{}{
			unit.Date,
			unit.Count,
		})
	}

	c.JSON(200, gin.H{
		"message":                 "获取个人信息成功",
		"recent_submit_statistic": recentSubmitStatistic,
	})
}
