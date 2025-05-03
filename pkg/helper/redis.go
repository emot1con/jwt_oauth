package helper

import (
	"auth/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func AddToBlacklist(token string) error {
	key := fmt.Sprintf("bl:%s", token)
	now := time.Now()

	if err := config.RedisClient.Set(context.Background(), key, "1", now.AddDate(0, 3, 0).Sub(now)).Err(); err != nil {
		logrus.Errorf("Error adding token to blacklist: %v", err)
		return err
	}

	return nil
}

func IsTokenBlacklisted(token string) (bool, error) {
	redisClient := config.RedisClient
	key := fmt.Sprintf("bl:%s", token)

	result, err := redisClient.Exists(context.Background(), key).Result()
	if err != nil {
		logrus.Errorf("Error checking token blacklist: %v", err)
		return false, err
	}

	return result > 0, nil
}
