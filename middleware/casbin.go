package middleware

import (
	"ahpuoj/model"
	"ahpuoj/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var sub string
		user, _ := c.Get("user")
		if user, ok := user.(model.User); ok {
			sub = user.Role
		}

		obj := c.Request.URL.RequestURI()
		act := c.Request.Method
		enforcer := model.GetCasbin()
		if res, err := enforcer.EnforceSafe(sub, obj, act); err != nil {
			utils.Consolelog(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "内部错误",
			})
			c.Abort()
			return
		} else if res {
			c.Next()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "很抱歉您没有此权限",
			})
			c.Abort()
		}
		c.Next()
	}
}
