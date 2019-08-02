package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	randN := random.Intn(200)
	file, _ := os.Create("/home/jiezi19971225/go/src/ahpuoj/test/test.out")
	os.Stdout = file
	fmt.Println(randN)
	fmt.Print(5, " ")
	for i := 1; i < randN; i++ {
		randV := random.Intn(100)
		if randV >= 50 {
			fmt.Print(5, " ")
		} else if randV >= 20 {
			fmt.Print(10, " ")
		} else {
			fmt.Print(20, " ")
		}
	}
}
