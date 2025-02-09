package main

import (
	"github.com/isaacmirandacampos/rate-limiter/configs"
)

func main() {
	configs, err := configs.LoadConfig(".env")
	if err != nil {
		panic(err)
	}
	println("Hello, World!")
	println(configs.Timeout)
}
