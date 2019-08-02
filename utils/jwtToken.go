package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var TokenSinature string

func init() {
	TokenSinature = "secret"
}

type MyCustomClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(username string) string {
	claims := &MyCustomClaims{
		username,
		jwt.StandardClaims{
			// 过期时间为2小时
			ExpiresAt: time.Now().Add(time.Hour * 2000).Unix(),
			// ExpiresAt: time.Now().Add(time.Second * 5).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(TokenSinature))
	return tokenString
}
