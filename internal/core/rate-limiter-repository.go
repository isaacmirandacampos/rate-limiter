package core

type RateLimiterRepository interface {
	Increment(key string, expiration int64) (int, error)
}
