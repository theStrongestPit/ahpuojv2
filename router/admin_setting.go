package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminSettingRouter(g *gin.RouterGroup) {
	g.GET("/settings", controller.GetSettings)
	g.PUT("/settings", controller.SetSettings)
}
