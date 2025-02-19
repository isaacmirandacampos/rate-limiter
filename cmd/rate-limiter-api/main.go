package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/isaacmirandacampos/rate-limiter/configs"
	"github.com/isaacmirandacampos/rate-limiter/internal/controller"
	"github.com/isaacmirandacampos/rate-limiter/internal/core"
	"github.com/isaacmirandacampos/rate-limiter/internal/database"
	"github.com/isaacmirandacampos/rate-limiter/internal/middleware"
)

func main() {
	configs, err := configs.LoadConfig(".env")
	if err != nil {
		panic(err)
	}
	timeout := time.Duration(configs.Timeout) * time.Second

	redisPool := database.NewRedisConnection(
		configs.RedisAddress,
	)
	defer redisPool.Close()
	redisRateLimiterRepository := core.NewRedisRateLimiterRepository(redisPool)
	rateLimiterHandler := core.NewRateLimiterHandler(redisRateLimiterRepository, int64(timeout.Seconds()))
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiterHandler, configs.RequestsPerSecondByIp, configs.RequestsPerSecondByApiKey)
	http.HandleFunc("/", rateLimiterMiddleware.Execute(controller.HelloWorld))
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
