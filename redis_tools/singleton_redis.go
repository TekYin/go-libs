package redis_tools

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client     *redis.Client
	once       sync.Once
	defaultTTL = 0 * time.Second // no TTL by default
)

// Init initializes the singleton Redis client.
// Call this once early in your CLI program.
func Init(host string, port int, pass string) {
	once.Do(func() {
		addr := fmt.Sprintf("%s:%d", host, port)
		if addr == "" {
			addr = "127.0.0.1:6379"
		}

		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
			DB:       0,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			log.Fatalf("failed to connect to Redis: %v", err)
		}
	})
}

// Set sets a value in a Redis hash (persistent dictionary)
func Set(ctx context.Context, dict string, key string, value interface{}) error {
	return client.HSet(ctx, dict, key, value).Err()
}

// Get gets a value from a Redis hash
func Get(ctx context.Context, dict, key string) (string, bool, error) {
	val, err := client.HGet(ctx, dict, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	return val, true, err
}

// Delete removes a key from a Redis hash
func Delete(ctx context.Context, dict, key string) error {
	return client.HDel(ctx, dict, key).Err()
}

// Keys returns all keys in a Redis hash
func Keys(ctx context.Context, dict string) ([]string, error) {
	return client.HKeys(ctx, dict).Result()
}

// All returns all key-value pairs from a Redis hash
func All(ctx context.Context, dict string) (map[string]string, error) {
	return client.HGetAll(ctx, dict).Result()
}

func Exists(ctx context.Context, dict, key string) (bool, error) {
	res, err := client.HExists(ctx, dict, key).Result()
	return res, err
}

func ExistsDict(ctx context.Context, dict string) (bool, error) {
	res, err := client.Exists(ctx, dict).Result()
	return res == 1, err
}
