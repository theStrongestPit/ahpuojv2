package utils

import (
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendMail(mailTo []string, subject string, body string) error {

	cfg := GetCfg()
	// cfg := GetTestCfg("../config/config.ini")

	mailcfg, _ := cfg.GetSection("mail")
	mailConn := map[string]string{
		"user": mailcfg["user"],
		"pass": mailcfg["pass"],
		"host": mailcfg["host"],
		"port": mailcfg["port"],
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int
	m := gomail.NewMessage()
	m.SetHeader("From", "XD Game"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo...)                            //发送给多个用户
	m.SetHeader("Subject", subject)                         //设置邮件主题
	m.SetBody("text/html", body)                            //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}
