package middleware

import (
	"fmt"
	"net/http"

	"github.com/isaacmirandacampos/rate-limiter/internal/core"
)

type RateLimiter struct {
	core.LimiterByIp
}

func NewRateLimiter(limiterByIp *core.LimiterByIp) *RateLimiter {
	return &RateLimiter{
		*limiterByIp,
	}
}

func (rate *RateLimiter) RateLimiterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr[:len(r.RemoteAddr)-6]
		header := r.Header.Get("API_KEY")
		fmt.Print("API_KEY: ", header)
		allow, err := rate.RateLimiterByIp(ip)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if !allow {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
