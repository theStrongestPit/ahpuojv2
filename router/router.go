package router

import (
	"ahpuoj/controller"
	"ahpuoj/middleware"

	"github.com/gin-gonic/gin"
)

func handelRouterGroup() {

}

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/login", controller.Login)
	router.POST("/api/register", controller.Register)
	router.POST("/api/findpass", controller.SendFindPassEmail)
	router.GET("/api/verifyresetpasstoken", controller.VeriryResetPassToken)
	router.POST("/api/resetpassbytoken", controller.ResetPassByToken)

	// 添加解析token中间件
	router.Use(middleware.ParseTokenMiddleware())
	// 无需用户登录的api
	nologin := router.Group("/api")
	ApiNologinRouter(nologin)
	// 添加JWT认证中间件
	router.Use(middleware.JwtauthMiddleware())
	// 需要用户登录的api
	user := router.Group("/api")
	ApiUserRouter(user)
	// 后台路由组 添加Casbin权限控制中间件
	admin := router.Group("/api/admin", middleware.CasbinMiddleware())

	// 权限控制路由
	ApiAdminCasbinRouter(admin)
	// 后台用户相关路由
	ApiAdminUserRouter(admin)
	// 后台用户问题集相关路由
	ApiAdminProblemsetRouter(admin)
	// 后台问题相关路由
	ApiAdminProblemRouter(admin)
	// 后台竞赛&作业相关路由
	ApiAdminContestRouter(admin)
	// 后台系列赛相关路由
	ApiAdminSeriesRouter(admin)
	// 后台标签相关路由
	ApiAdminTagRouter(admin)
	// 后台用户生成器相关路由
	ApiAdminGeneratorRouter(admin)
	// 后台团队相关路由
	ApiAdminTeamRouter(admin)
	// 后台新闻相关路由
	ApiAdminNewRouter(admin)
	// 后台图片相关路由
	ApiAdminImageRouter(admin)
	// 后台讨论相关路由
	ApiAdminIssueRouter(admin)
	// 后台配置项相关路由
	ApiAdminSettingRouter(admin)
	// 后台其他路由
	ApiAdminRouter(admin)

	return router
}
