package test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMask(t *testing.T) {
	langSupport := []bool{
		true,  //0  c
		true,  //1  cpp
		false, //2 pascal
		true,  //3  java
		false, //4 ruby
		false, //5 bash
		true,  //6  python
		false, //7 php
		false, //8 perl
		false, //9 csharp
		false, //10    ojb-c
		false, //11    freebasic
		false, //12    scheme
		false, //13    clang
		false, //14    clang++
		false, //15    lua
		false, //16    js
		true,  //17 golang
	}
	mask := 0

	for i := 0; i < len(langSupport); i++ {
		if langSupport[i] {
			mask = mask | (1 << uint(i))
		}
	}
	t.Log(mask)
}
