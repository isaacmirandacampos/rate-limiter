package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/isaacmirandacampos/rate-limiter/internal/core"
	"github.com/isaacmirandacampos/rate-limiter/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func setupRedisTestContainer(t *testing.T) (*redigo.Pool, func()) {
	ctx := context.Background()
	redisContainer, err := redis.Run(ctx,
		"redis:7.4.2-alpine",
	)
	assert.NoError(t, err)

	redisHost, err := redisContainer.Host(ctx)
	assert.NoError(t, err)
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	assert.NoError(t, err)
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort.Port())

	redisPool := &redigo.Pool{
		MaxIdle:     10,
		MaxActive:   50,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", redisAddr)
		},
	}

	cleanup := func() {
		redisContainer.Terminate(ctx)
	}

	return redisPool, cleanup
}

func testRateLimiter(t *testing.T, rateLimiterMiddleware *middleware.RateLimiterMiddleware, headers map[string]string, expectedPasses int) {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimiterMiddleware.Execute(mockHandler)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	for i := 1; i <= expectedPasses; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Request %d should pass", i)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request %d should be rate limited", expectedPasses+1)
}

func TestRateLimiterMiddlewareByIp(t *testing.T) {
	redisPool, cleanup := setupRedisTestContainer(t)
	defer cleanup()

	repo := core.NewRedisRateLimiterRepository(redisPool)
	rateLimiter := middleware.NewRateLimiterMiddleware(core.NewRateLimiterHandler(repo, 3), 2, 1)

	testRateLimiter(t, rateLimiter, nil, 2)
}

func TestRateLimiterMiddlewareByApiKey(t *testing.T) {
	redisPool, cleanup := setupRedisTestContainer(t)
	defer cleanup()

	repo := core.NewRedisRateLimiterRepository(redisPool)
	rateLimiter := middleware.NewRateLimiterMiddleware(core.NewRateLimiterHandler(repo, 3), 1, 2)

	testRateLimiter(t, rateLimiter, map[string]string{"API_KEY": "1234"}, 2)
}
