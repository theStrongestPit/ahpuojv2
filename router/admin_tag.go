package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminTagRouter(g *gin.RouterGroup) {
	g.POST("/tag", controller.StoreTag)
	g.GET("/tags", controller.IndexTag)
	g.GET("/alltags", controller.GetAllTags)
	g.PUT("/tag/:id", controller.UpdateTag)
	g.DELETE("/tag/:id", controller.DeleteTag)
}
