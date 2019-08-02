package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyCustomClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

func createToken(username string) string {
	claims := MyCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2000).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("secret"))
	return tokenString
}

func main() {
	fmt.Println(createToken("admin"))
}
