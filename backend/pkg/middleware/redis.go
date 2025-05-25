package middleware

import (
	"auth/internal/config"
	"context"
	"time"
)

func RateLimiter(ctx context.Context, key string, limit int, duration time.Duration) (bool, error) {
	count, err := config.RedisClient.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		config.RedisClient.Expire(ctx, key, duration)
	}

	if count > int64(limit) {
		return false, nil
	}

	return true, nil
}
