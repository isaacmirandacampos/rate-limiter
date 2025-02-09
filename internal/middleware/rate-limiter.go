package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RateLimiter struct {
	Timeout           time.Duration
	requestsPerSecond int32
	*redis.Pool
}

func NewRateLimiter(pool *redis.Pool, requestsPerSecond int32, timeout time.Duration) *RateLimiter {
	return &RateLimiter{
		Timeout:           timeout,
		requestsPerSecond: requestsPerSecond,
		Pool:              pool,
	}
}

func (rate *RateLimiter) RateLimiterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn := rate.Pool.Get()
		defer conn.Close()
		ip := r.RemoteAddr[:len(r.RemoteAddr)-6]
		redisQ := `
			local current = redis.call("INCR", KEYS[1])
			if current == 1 then
				redis.call("EXPIRE", KEYS[1], ARGV[1])
			end
			return current
		`
		limit, err := redis.Int(conn.Do("EVAL", redisQ, 1, ip, int(rate.Timeout.Seconds())))
		if err != nil {
			fmt.Printf("error on Redis: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Requests from %s: %d/%d\n", ip, limit, rate.requestsPerSecond)
		if int32(limit) > rate.requestsPerSecond {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		reply, err := conn.Do("INCR", ip, "EX")
		if err != nil {
			fmt.Printf("error on incr: %v \n", err)
		}
		fmt.Printf("actual value of redis: %v \n", reply)
		conn.Flush()
		next(w, r)
	}
}
