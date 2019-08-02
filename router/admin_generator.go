package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminGeneratorRouter(g *gin.RouterGroup) {
	g.POST("/generator/compete", controller.CompeteAccountGenerator)
	g.POST("/generator/user", controller.UserAccountGenerator)
}
