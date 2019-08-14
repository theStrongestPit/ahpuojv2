package controller

import (
	"ahpuoj/model"
	"ahpuoj/request"
	"ahpuoj/utils"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req request.Login
	var user model.User
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}

	err = DB.Get(&user, "select * from user where username = ?", req.Username)

	if utils.CheckError(c, err, "用户不存在") != nil {
		return
	}

	h := sha1.New()
	h.Write([]byte(user.PassSalt))
	h.Write([]byte(req.Password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))

	if hashedPassword != user.Password {
		c.JSON(400, gin.H{
			"message": "密码错误",
		})
	} else {
		// 根据用户名payload生成token
		token := utils.CreateToken(req.Username)

		// 更新redis的token,过期时间为15天
		utils.Consolelog("登录成功")
		utils.Consolelog(token)
		conn := REDISPOOL.Get()
		defer conn.Close()
		reply, err := conn.Do("set", "token:"+req.Username, token)
		utils.Consolelog(reply, err)
		conn.Do("expire", "token:"+req.Username, 60*60*24*15)

		// 设置cookies
		cfg := utils.GetCfg()
		domain, _ := cfg.GetValue("project", "server")
		cookieLiveTimeStr, _ := cfg.GetValue("project", "cookielivetime")
		cookieLiveTime, _ := strconv.Atoi(cookieLiveTimeStr)
		c.SetCookie("access-token", token, cookieLiveTime, "/", domain, false, false)
		c.JSON(200, gin.H{
			"message": "登录成功",
		})
	}
}

func Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}

	// 加盐处理 16位随机字符串
	salt := utils.GetRandomString(16)

	h := sha1.New()
	h.Write([]byte(salt))
	h.Write([]byte(req.Password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))

	user := model.User{
		Username: req.Username,
		Nick:     req.Nick,
		Email:    sql.NullString{req.Email, true},
		Password: hashedPassword,
		PassSalt: salt,
	}

	err = user.Save()

	if utils.CheckError(c, err, "该邮箱/用户名/昵称已被注册") != nil {
		return
	}

	token := utils.CreateToken(user.Username)

	// 更新redis的token,过期时间为15天
	conn := REDISPOOL.Get()
	defer conn.Close()
	conn.Do("set", "token:"+user.Username, token)
	conn.Do("expire", "token:"+user.Username, 60*60*24*15)

	// 设置cookies
	cfg := utils.GetCfg()
	domain, _ := cfg.GetValue("project", "server")
	cookieLiveTimeStr, _ := cfg.GetValue("project", "cookielivetime")
	cookieLiveTime, _ := strconv.Atoi(cookieLiveTimeStr)
	c.SetCookie("access-token", token, cookieLiveTime, "/", domain, false, false)
	c.JSON(200, gin.H{
		"message": "注册成功",
		"token":   token,
	})
}

// 发送重设密码邮件的接口
func SendFindPassEmail(c *gin.Context) {
	var req request.ResetPass
	err := c.ShouldBindJSON(&req)

	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}

	var user_id int
	err = DB.Get(&user_id, "select id from user where email = ?", req.Email)

	if utils.CheckError(c, err, "用户不存在") != nil {
		return
	}

	// 生成随机字符串
	token := utils.GetRandomString(30)
	_, err = DB.Exec("insert into resetpassword(user_id,token,expired_at) values(?,?,date_add(NOW(),INTERVAL 1 hour)) on duplicate key update token = ?,expired_at=date_add(NOW(),INTERVAL 1 hour)", user_id, token, token)
	cfg := utils.GetCfg()
	server, _ := cfg.GetValue("project", "server")

	mailTo := []string{
		req.Email,
	}
	//邮件主题
	subject := "AHPUOJ重设密码邮件"
	// 邮件正文
	body := fmt.Sprintf("请访问以下连接重设您的密码，链接将会在1小时内失效，请尽快进行设置 <a href=\"%s/resetpassword?token=%s\">%s/resetpass?token=%s</a>", server, token, server, token)
	utils.SendMail(mailTo, subject, body)
	c.JSON(200, gin.H{
		"message": "已成功发送重设密码邮件，请前往邮箱查看",
	})
}

func VeriryResetPassToken(c *gin.Context) {
	token := c.Query("token")
	type T struct {
		Token      string `db:"token"`
		Expired_at string `db:"expired_at"`
	}
	var t T
	err := DB.Get(&t, "select token,expired_at from  resetpassword where token = ?", token)

	if utils.CheckError(c, err, "token非法") != nil {
		return
	}

	now := time.Now()
	expire, _ := time.Parse("2006-01-02 15:04:05", t.Expired_at)
	if now.After(expire) {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "token已过期，请重新发送邮件",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "token验证成功，请立即重设密码",
	})
}
func ResetPassByToken(c *gin.Context) {
	var req request.ResetPassByToken
	err := c.ShouldBindJSON(&req)

	if utils.CheckError(c, err, "参数错误") != nil {
		return
	}

	type T struct {
		User_id    string `db:"user_id"`
		Token      string `db:"token"`
		Expired_at string `db:"expired_at"`
	}
	// 验证token
	var t T
	err = DB.Get(&t, "select user_id,token,expired_at from  resetpassword where token = ?", req.Token)

	if utils.CheckError(c, err, "token非法") != nil {
		return
	}

	now := time.Now()
	expire, _ := time.Parse("2006-01-02 15:04:05", t.Expired_at)
	if now.After(expire) {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "token已过期，请重新发送邮件",
		})
		return
	}

	// 更新密码
	// 加盐处理 16位随机字符串
	h := sha1.New()
	salt := utils.GetRandomString(16)
	h.Reset()
	h.Write([]byte(salt))
	h.Write([]byte(req.Password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))

	_, err = DB.Exec("update user set password = ?, passsalt = ? where id = ?", hashedPassword, salt, t.User_id)
	_, err = DB.Exec("delete from resetpassword where token = ?", t.Token)
	c.JSON(200, gin.H{
		"message": "密码修改成功",
	})

}
