package utils

import (
	"testing"
)

func TestMain(t *testing.T) {
	mailTo := []string{
		"1755465135@qq.com",
	}
	//邮件主题为"Hello"  
	subject := "Hello"
	// 邮件正文
	body := "Good"
	// t.Log(mailTo, subject, body)
	SendMail(mailTo, subject, body)
}
