package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminUserRouter(g *gin.RouterGroup) {
	g.GET("/users", controller.IndexUser)
	g.PUT("/user/:id/status", controller.ToggleUserStatus)
	g.PUT("/user/:id/pass", controller.ChangeUserPass)
	// g.PUT("/user/:id", controller.UpdateUser)
}
