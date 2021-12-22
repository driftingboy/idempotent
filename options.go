package idempotent

import (
	"time"

	redis "github.com/go-redis/redis/v7"
)

// Required
type Config struct {
	RedisAddrs []string
	Password   string
}

// Not required
type Option func(*redis.ClusterOptions)

func WithTimeout(dial, read, write time.Duration) Option {
	return func(o *redis.ClusterOptions) {
		o.DialTimeout = dial
		o.ReadTimeout = read
		o.WriteTimeout = write
	}
}
