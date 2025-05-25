package helper

import (
	"auth/internal/config"
	"context"
	"fmt"
	"time"
)

func RateLimiter(ctx context.Context, key string) (bool, error) {
	redisKey := fmt.Sprintf("rate_limiter:%s", key)

	count, err := config.RedisClient.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := config.RedisClient.Expire(ctx, redisKey, 1*time.Minute).Err(); err != nil {
			return false, err
		}
	}

	if count > 7 {
		return false, nil
	}

	return true, nil
}
