package core

import (
	"errors"
)

type RateLimiterHandler struct {
	rateLimiterRepository RateLimiterRepository
	timeout               int64
}

func NewRateLimiterByIp(rateLimiterRepository RateLimiterRepository, timeout int64) *RateLimiterHandler {
	return &RateLimiterHandler{
		rateLimiterRepository: rateLimiterRepository,
		timeout:               timeout,
	}
}

func (l *RateLimiterHandler) Execute(key string, requestsPerSecond int32) (bool, error) {
	limit, err := l.rateLimiterRepository.Increment(key, l.timeout)
	if err != nil {
		return false, errors.New("internal server error")
	}
	if int32(limit) > requestsPerSecond {
		return false, nil
	}
	return true, nil
}
