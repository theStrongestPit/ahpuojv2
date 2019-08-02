package middleware

import (
	"ahpuoj/controller"
	"ahpuoj/model"
	"ahpuoj/utils"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func parseToken(c *gin.Context) (model.User, error) {

	tokenString := c.GetHeader("Authorization")
	var user model.User
	var role string

	token, err := jwt.ParseWithClaims(tokenString, &utils.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.TokenSinature), nil
	})

	// 忽略超时错误
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&ve.Errors&jwt.ValidationErrorExpired != 0 {
		} else {
			return user, err
		}
	}

	if claims, ok := token.Claims.(*utils.MyCustomClaims); ok {
		c.Set("tokenExpireAt", claims.ExpiresAt)
		username := claims.UserName
		err = model.DB.Get(&user, "select * from user where username = ?", username)
		if err != nil {
			return user, errors.New("用户不存在")
		}
		err = model.DB.Get(&role, "select role.name from role where id = ?", user.RoleId)
		utils.Consolelog("role", role)
		user.Role = role
		// 判断用户登录token是否存在redis缓存中
		conn := controller.REDISPOOL.Get()
		defer conn.Close()
		storeToken, _ := redis.String(conn.Do("get", "token:"+username))

		if storeToken != tokenString {
			return user, errors.New("token已被废弃")
		}
	} else {
		return user, errors.New("token结构不匹配")
	}
	return user, nil
}

func ParseTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if token != nil {
			// 判断用户是否存在
			user, err := parseToken(c)
			if err != nil {
				c.Next()
			} else {
				c.Set("user", user)
				c.Next()
			}
		}
	}
}
