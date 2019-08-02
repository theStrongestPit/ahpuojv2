package router

import (
	"ahpuoj/controller"

	"github.com/gin-gonic/gin"
)

func ApiAdminSeriesRouter(g *gin.RouterGroup) {
	g.POST("/series", controller.StoreSeries)
	g.GET("/serieses", controller.IndexSeries)
	g.GET("/series/:id", controller.GetSeries)
	g.GET("/series/:id/contests", controller.IndexSeriesContest)
	g.PUT("/series/:id", controller.UpdateSeries)
	g.DELETE("/series/:id", controller.DeleteSeries)
	g.PUT("/series/:id/status", controller.ToggleSeriesStatus)
	g.POST("/series/:id/contest/:contestid", controller.AddSeriesContest)
	g.DELETE("/series/:id/contest/:contestid", controller.DeleteSeriesContest)
}
