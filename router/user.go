package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiUserRouter(g *gin.RouterGroup) {
	g.GET("/user", controller.GetUser)
	g.POST("/submit", controller.SubmitToJudge)
	g.POST("/testrun", controller.SubmitToTestRun)
	g.POST("/issue", controller.PostIssue)
	g.POST("/issue/:id/reply", controller.ReplyToIssue)
	g.GET("/myreplys", controller.GetMyReplys)
	g.PUT("/solution/:id/status", controller.ToggleSolutionStatus)
	g.PUT("/user/avatar", controller.UploadAvatar)
	g.PUT("/user/nick", controller.ResetNick)
	g.PUT("/user/password", controller.ResetPassword)
	g.GET("/datafile", controller.DownloadDataFile)
	g.GET("/problem/:id/latestsource", controller.GetLatestSource)
	g.GET("/contest/:id/problem/:num/latestsource", controller.GetLatestContestSource)
}
