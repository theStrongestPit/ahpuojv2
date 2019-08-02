package main

import (
	"ahpuoj/router"
)

func main() {
	r := router.InitRouter()
	// Listen and Servesr in 0.0.0.0:8080
	r.Run(":8080")
}
