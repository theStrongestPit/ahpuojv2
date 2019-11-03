package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminProblemRouter(g *gin.RouterGroup) {
	g.POST("/problem", controller.StoreProblem)
	g.GET("/problems", controller.IndexProblem)
	g.GET("/problem/:id", controller.GetProblem)
	g.PUT("/problem/:id", controller.UpdateProblem)
	g.DELETE("/problem/:id", controller.DeleteProblem)
	g.PUT("/problem/:id/status", controller.ToggleProblemStatus)

	g.GET("/problem/:id/datas", controller.IndexProblemData)
	g.GET("/problem/:id/data/:filename", controller.GetProblemData)
	g.POST("/problem/:id/data", controller.AddProblemData)
	g.POST("/problem/:id/datafile", controller.AddProblemDataFile)
	g.PUT("/problem/:id/data/:filename", controller.EditProblemData)
	g.DELETE("/problem/:id/data/:filename", controller.DeleteProblemData)

	// 重判问题路由
	g.PUT("/solution/:id/judgestatus", controller.RejudgeSolution)
	g.PUT("/problem/:id/judgestatus", controller.RejudgeProblem)
	// 重排问题路由
	g.PUT("/problem/:id/movement/:newid", controller.ReassignProblem)
}
