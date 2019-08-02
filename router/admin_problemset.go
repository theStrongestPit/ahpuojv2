package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminProblemsetRouter(g *gin.RouterGroup) {
	g.POST("/problemset", controller.ImportProblemSet)
}
