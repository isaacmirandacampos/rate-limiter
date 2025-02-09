package main

import (
	"fmt"
	"net/http"

	"github.com/isaacmirandacampos/rate-limiter/configs"
)

func main() {
	configs, err := configs.LoadConfig(".env")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!\n")
		fmt.Fprintf(w, "Timeout %v", configs.Timeout)
	})
	http.ListenAndServe(":8080", nil)
}
