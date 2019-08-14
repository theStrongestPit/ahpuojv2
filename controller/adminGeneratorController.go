package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CompeteAccountGenerator(c *gin.Context) {
	var req request.CompeteAccount
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}

	var infos []string
	var users []map[string]interface{}
	for i := 1; i <= req.Number; i++ {

		username := req.Prefix + strconv.Itoa(i)
		randomPassword := utils.GetRandomString(15)

		h := sha1.New()
		h.Write([]byte(randomPassword))
		hashedPassword := fmt.Sprintf("%x", h.Sum(nil))

		user := model.User{
			Username:      username,
			Nick:          username,
			Email:         sql.NullString{"", true},
			Password:      hashedPassword,
			IsCompeteUser: 1,
		}

		err = user.Save()
		if err != nil {
			infos = append(infos, "用户"+username+"创建失败")
		} else {
			users = append(users, map[string]interface{}{
				"username": username,
				"password": randomPassword,
			})

			infos = append(infos, "用户"+username+"创建成功")
		}

	}
	utils.Consolelog(users)
	c.JSON(200, gin.H{
		"message": "操作成功",
		"users":   users,
		"info":    infos,
	})
}

func UserAccountGenerator(c *gin.Context) {
	var req request.UserAccount
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}

	pieces := strings.Split(req.UserList, "\n")

	var infos []string
	var users []interface{}
	if len(pieces) > 0 && len(pieces[0]) > 0 {
		for _, username := range pieces {

			password := "123456"
			// 更新密码
			// 加盐处理 16位随机字符串
			salt := utils.GetRandomString(16)
			h := sha1.New()
			h.Write([]byte(salt))
			h.Write([]byte(password))
			hashedPassword := fmt.Sprintf("%x", h.Sum(nil))

			user := model.User{
				Username: username,
				Nick:     username,
				Email:    sql.NullString{"", true},
				Password: hashedPassword,
				PassSalt: salt,
			}
			err = user.Save()
			if err == nil {
				users = append(users, map[string]interface{}{
					"username": username,
					"password": password,
				})
				infos = append(infos, "用户"+username+"创建成功")
			} else {
				infos = append(infos, "用户"+username+"创建失败")
			}

		}
	}
	c.JSON(200, gin.H{
		"message": "操作成功",
		"users":   users,
		"info":    infos,
	})
}
