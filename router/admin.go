package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminRouter(g *gin.RouterGroup) {
	g.GET("/submitstatistic", controller.GetSubmitStatistic)
}
