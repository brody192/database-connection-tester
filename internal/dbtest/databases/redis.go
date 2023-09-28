package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func Redis(redisURL string) (time.Duration, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return 0, fmt.Errorf("redis redis.ParseURL error: %w", err)
	}

	opt.DialTimeout = 10 * time.Second
	opt.MaxRetries = 500
	opt.MinRetryBackoff = 10 * time.Millisecond

	sT := time.Now()

	client := redis.NewClient(opt)

	defer client.Close()

	if err := client.Ping(context.Background()).Err(); err != nil {
		return time.Since(sT), fmt.Errorf("client.Ping error: %w", err)
	}

	return time.Since(sT), nil
}
