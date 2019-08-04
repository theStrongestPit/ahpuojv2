package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminIssueRouter(g *gin.RouterGroup) {
	g.PUT("/issue/:id/status", controller.ToggleIssueStatus)
	g.PUT("/reply/:id/status", controller.ToggleReplyStatus)
}
