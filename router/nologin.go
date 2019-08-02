package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

// 游客路由
func ApiNologinRouter(g *gin.RouterGroup) {
	g.GET("/news", controller.NologinGetNewList)
	g.GET("/problem/:id", controller.NologinGetProblem)
	g.GET("/problems", controller.NologinGetProblemList)
	g.GET("/contest/:id", controller.NologinGetContest)
	g.GET("/contest/:id/ranklist", controller.NologinGetContestRankList)
	g.GET("/contest/:id/teamranklist", controller.NologinGetContestTeamRankList)
	g.GET("/contest/:id/problem/:num", controller.NologinGetContestProblem)
	g.GET("/contests", controller.NologinGetContestList)
	g.GET("/solution/:id", controller.NologinGetSolution)
	g.GET("/solutions", controller.NologinGetSolutionList)
	g.GET("/alltags", controller.NologinGetAllTags)
	g.GET("/languages", controller.NologinGetLanguageList)
	g.GET("/problem/:id/issues", controller.NologinGetIssueList)
	g.GET("/issue/:id", controller.NologinGetIssue)
	g.GET("/user/:id", controller.NologinGetUserInfo)
	g.GET("/ranklist", controller.NologinGetRankList)
}
