package middleware

import (
	"net/http"

	"github.com/isaacmirandacampos/rate-limiter/internal/core"
)

type RateLimiter struct {
	core.RateLimiterHandler
	requestsPerSecondByIp     int32
	requestsPerSecondByApiKey int32
}

func NewRateLimiter(rateLimiterHandler *core.RateLimiterHandler, requestsPerSecondByIp int32, requestsPerSecondByApiKey int32) *RateLimiter {
	return &RateLimiter{
		RateLimiterHandler:        *rateLimiterHandler,
		requestsPerSecondByIp:     requestsPerSecondByIp,
		requestsPerSecondByApiKey: requestsPerSecondByApiKey,
	}
}

func (rate *RateLimiter) RateLimiterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr[:len(r.RemoteAddr)-6]
		apiKey := r.Header.Get("API_KEY")
		var allow bool
		var err error
		if apiKey != "" {
			allow, err = rate.Execute(apiKey, rate.requestsPerSecondByApiKey)
		} else {
			allow, err = rate.Execute(ip, rate.requestsPerSecondByIp)
		}
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
