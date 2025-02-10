package core

import (
	"errors"
)

type LimiterByIp struct {
	rateLimiterRepository RateLimiterRepository
	requestsPerSecond     int32
	timeout               int64
}

func NewRateLimiterByIp(rateLimiterRepository RateLimiterRepository, requestsPerSecond int32, timeout int64) *LimiterByIp {
	return &LimiterByIp{
		rateLimiterRepository: rateLimiterRepository,
		requestsPerSecond:     requestsPerSecond,
		timeout:               timeout,
	}
}

func (l *LimiterByIp) RateLimiterByIp(ip string) (bool, error) {
	limit, err := l.rateLimiterRepository.Increment(ip, l.timeout)
	if err != nil {
		return false, errors.New("internal server error")
	}
	if int32(limit) > l.requestsPerSecond {
		return false, nil
	}
	return true, nil
}
