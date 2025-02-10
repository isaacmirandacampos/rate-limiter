package core

import (
	"github.com/gomodule/redigo/redis"
)

type RedisRateLimiterRepository struct {
	redisPool *redis.Pool
}

func NewRedisRateLimiterRepository(pool *redis.Pool) *RedisRateLimiterRepository {
	return &RedisRateLimiterRepository{
		redisPool: pool,
	}
}

func (r *RedisRateLimiterRepository) Increment(key string, expiration int64) (int, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	redisQ := `
			local current = redis.call("INCR", KEYS[1])
			if current == 1 then
				redis.call("EXPIRE", KEYS[1], ARGV[1])
			end
			return current
		`
	limit, err := redis.Int(conn.Do("EVAL", redisQ, 1, key, expiration))
	if err != nil {
		return 0, err
	}
	return limit, nil
}
