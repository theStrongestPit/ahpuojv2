package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminTeamRouter(g *gin.RouterGroup) {
	g.POST("/team", controller.StoreTeam)
	g.POST("/team/:id/users", controller.AddTeamUsers)
	g.GET("/teams", controller.IndexTeam)
	// 在竞赛作业人员配置中获取全部团队的接口
	g.GET("/allteams", controller.GetAllTeams)
	g.GET("/team/:id", controller.GetTeam)
	g.GET("/team/:id/users", controller.IndexTeamUser)
	g.PUT("/team/:id", controller.UpdateTeam)
	g.DELETE("/team/:id/user/:userid", controller.DeleteTeamUser)
	g.DELETE("/team/:id", controller.DeleteTeam)
}
