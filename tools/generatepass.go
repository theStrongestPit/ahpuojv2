package main

import (
	"crypto/sha1"
	"fmt"
)

func main() {
	salt := "abcdefghijklmnop"
	h := sha1.New()
	h.Write([]byte(salt))
	h.Write([]byte("123456"))
	fmt.Print(fmt.Sprintf("%x", h.Sum(nil)))
}
