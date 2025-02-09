package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RateLimiter struct {
	Timeout time.Duration
	*redis.Pool
}

func NewRateLimiter(pool *redis.Pool, timeout time.Duration) *RateLimiter {
	return &RateLimiter{
		Timeout: timeout,
		Pool:    pool,
	}
}

func (rate *RateLimiter) RateLimiterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn := rate.Pool.Get()
		defer conn.Close()
		ip := r.RemoteAddr[:len(r.RemoteAddr)-6]
		reply, err := conn.Do("INCR", ip)
		if err != nil {
			fmt.Printf("error on incr: %v \n", err)
		}
		fmt.Printf("actual value of redis: %v \n", reply)
		conn.Flush()
		next(w, r)
	}
}
