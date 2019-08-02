package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminContestRouter(g *gin.RouterGroup) {
	g.POST("/contest", controller.StoreContest)
	g.GET("/contests", controller.IndexContest)
	g.GET("/allcontests", controller.GetAllContests)
	g.GET("/contest/:id", controller.GetContest)
	g.PUT("/contest/:id", controller.UpdateContest)
	g.DELETE("/contest/:id", controller.DeleteContest)
	g.PUT("/contest/:id/status", controller.ToggleContestStatus)
	g.GET("/contest/:id/users", controller.IndexContestUser)
	g.POST("/contest/:id/users", controller.AddContestUsers)
	g.DELETE("/contest/:id/user/:userid", controller.DeleteContestUser)
	g.GET("/contest/:id/teams", controller.IndexContestTeamWithUser)
	g.POST("/contest/:id/team/:teamid", controller.AddContestTeam)
	g.DELETE("/contest/:id/team/:teamid", controller.DeleteContestTeam)
	g.POST("/contest/:id/team/:teamid/users", controller.AddContestTeamUsers)
	g.POST("/contest/:id/team/:teamid/allusers", controller.AddContestTeamAllUsers)
	g.DELETE("/contest/:id/team/:teamid/user/:userid", controller.DeleteContestTeamUser)
}
