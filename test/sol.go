package main

import (
	"fmt"
	"os"
)

func main() {
	var n int
	inFile, _ := os.Open("/home/jiezi19971225/go/src/ahpuoj/test/test.out")
	os.Stdin = inFile
	fmt.Scanf("%d", &n)
	fmt.Println(n)
	five, ten := 0, 0
	for i := 0; i < n; i++ {
		var v int
		fmt.Scanf("%v", &v)
		fmt.Println(i, v)
		if v == 5 {
			five++
		}
		if v == 10 {
			if five >= 1 {
				five--
				ten++
			} else {
				fmt.Println("NO")
				return
			}
		}
		if v == 20 {
			if ten >= 1 && five >= 1 {
				five--
				ten--
			} else if five >= 3 {
				five -= 3
			} else {
				fmt.Println("NO")
				return
			}
		}
	}
	fmt.Println("OK")
}
