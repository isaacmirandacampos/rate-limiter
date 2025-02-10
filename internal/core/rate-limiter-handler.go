package core

import (
	"errors"
)

type LimiterByIp struct {
	rateLimiterRepository RateLimiterRepository
	timeout               int64
}

func NewRateLimiterByIp(rateLimiterRepository RateLimiterRepository, requestsPerSecond int32, timeout int64) *LimiterByIp {
	return &LimiterByIp{
		rateLimiterRepository: rateLimiterRepository,
		timeout:               timeout,
	}
}

func (l *LimiterByIp) RateLimiterHandler(key string, requestsPerSecond int32) (bool, error) {
	limit, err := l.rateLimiterRepository.Increment(key, l.timeout)
	if err != nil {
		return false, errors.New("internal server error")
	}
	if int32(limit) > requestsPerSecond {
		return false, nil
	}
	return true, nil
}
