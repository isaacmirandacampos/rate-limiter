package database

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func NewRedisConnection(
	address string,
) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   50,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", address)
			if err != nil {
				log.Fatalf("Failed to connect to Redis: %v", err)
			}
			return conn, err
		},
	}
}
