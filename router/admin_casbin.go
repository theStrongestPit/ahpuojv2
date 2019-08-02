package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminCasbinRouter(g *gin.RouterGroup) {
	g.POST("/casbin", controller.StoreCasbin)
}
