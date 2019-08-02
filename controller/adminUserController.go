package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexUser(c *gin.Context) {

	pageStr := c.Query("page")
	perpageStr := c.Query("perpage")
	userType := c.Query("userType")
	param := c.Query("param")
	page, _ := strconv.Atoi(pageStr)
	perpage, _ := strconv.Atoi(perpageStr)
	if page == 0 {
		page = 1
	}
	whereString := " where is_compete_user =" + userType
	if len(param) > 0 {
		whereString += " and username like '%" + param + "%' or nick like '%" + param + "%'"
	}
	whereString += " order by id desc"
	rows, total, err := model.Paginate(page, perpage, "user", []string{"*"}, whereString)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "数据获取失败",
		})
		return
	}
	var users []map[string]interface{}
	for rows.Next() {
		var user model.User
		err = rows.StructScan(&user)
		if err != nil {
			utils.Consolelog(err)
		}
		users = append(users, user.Response())
	}
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"total":   total,
		"perpage": perpage,
		"data":    users,
	})
}

func ToggleUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := model.User{
		Id: id,
	}

	err := user.ToggleStatus()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "更改用户状态失败，用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改用户状态成功",
	})
}

func ChangeUserPass(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UserPass
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "请求参数错误",
		})
		return
	}

	// 更新密码
	// 加盐处理 16位随机字符串
	salt := utils.GetRandomString(16)
	h := sha1.New()
	h.Write([]byte(salt))
	h.Write([]byte(req.Password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))

	user := model.User{
		Id:       id,
		Password: hashedPassword,
		PassSalt: salt,
	}

	err = user.ChangePass()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "更改用户密码失败，用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更改用户密码成功",
	})
}
