package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminImageRouter(g *gin.RouterGroup) {
	g.POST("/image", controller.StoreImage)
}
