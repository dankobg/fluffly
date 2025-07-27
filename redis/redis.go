package redis

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/dankobg/fluffly/config"
	"github.com/redis/go-redis/v9"
)

// New returns the new [redis.Client]
func New(config config.RedisConfig) (*redis.Client, error) {
	addr := net.JoinHostPort(config.Host, strconv.Itoa(config.Port))

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis server: %w", err)
	}

	return rdb, nil
}
