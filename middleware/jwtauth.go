package middleware

import (
	"ahpuoj/controller"
	"ahpuoj/model"
	"ahpuoj/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JwtauthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if token != nil && token.Valid {
			// 判断用户是否存在
			if _, exist := c.Get("user"); exist {
				user, _ := controller.GetUserInstance(c)
				if user.Defunct == 0 {
					c.Next()
				} else {
					c.AbortWithStatusJSON(400, gin.H{
						"message": "该账号已被管理员封禁，请联系管理员解封",
					})
				}
			} else {
				c.AbortWithStatusJSON(400, gin.H{
					"message": "认证失败,请重新登录",
				})
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				c.AbortWithStatusJSON(400, gin.H{
					"message": "token非法",
				})
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				fmt.Println("token超时")
				// 判断用户是否存在 如果过期时间在两星期以内 生成新的token并返回
				if _, exist := c.Get("user"); exist {
					expire_at, _ := c.Get("tokenExpireAt")
					fmt.Println(expire_at)
					// 过期时间加上两星期 大于当前时间 可以刷新
					if expire_at, ok := expire_at.(int64); ok {
						twoWeeks := int64(3600 * (24*15 - 2))
						utils.Consolelog(expire_at, twoWeeks, time.Now().Unix())
						if expire_at+twoWeeks > time.Now().Unix() {
							fmt.Println("token在刷新时间内")
							user, _ := c.Get("user")
							if user, ok := user.(model.User); ok {
								newToken := utils.CreateToken(user.Username)
								conn := controller.REDISPOOL.Get()
								defer conn.Close()
								conn.Do("set", "token:"+user.Username, newToken)
								conn.Do("expire", "token:"+user.Username, 60*60*24*15)
								// 返回新的token
								cfg := utils.GetCfg()
								domain, _ := cfg.GetValue("project", "server")
								cookieLiveTimeStr, _ := cfg.GetValue("project", "cookielivetime")
								cookieLiveTime, _ := strconv.Atoi(cookieLiveTimeStr)
								fmt.Println("返回新的token")
								c.SetCookie("access-token", newToken, cookieLiveTime, "/", domain, false, false)
							}
						}
						c.Next()
					} else {
						c.AbortWithStatusJSON(400, gin.H{
							"message": "token已过期，请重新登录",
						})
					}
				} else {
					c.AbortWithStatusJSON(400, gin.H{
						"message": "token认证失败",
					})
				}
			} else {
				c.AbortWithStatusJSON(400, gin.H{
					"message": "无法处理token",
				})
			}
		} else {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "token认证失败",
			})
		}
	}
}
