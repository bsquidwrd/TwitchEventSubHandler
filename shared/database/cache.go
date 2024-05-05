package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	db *redis.Client
}

func NewCacheService() *CacheService {
	cacheUrl := os.Getenv("CACHE_URL")
	opt, err := redis.ParseURL(cacheUrl)
	if err != nil {
		panic(err)
	}
	slog.Debug("Cache connected successfully")

	client := redis.NewClient(opt)
	_, err = client.Get(context.Background(), "thiswillmostlikelyneverexistandthatsokay").Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	return &CacheService{
		db: client,
	}
}

func (c *CacheService) Cleanup() {
	defer c.db.Close()
}

func (c *CacheService) Ping() error {
	err := c.db.Echo(context.Background(), "OK").Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheService) TakeLock(key string, value string, expiration time.Duration) bool {
	lockKey := fmt.Sprintf("lock:%s", key)
	// Have a reasonable default a lock can be taken for
	if expiration > (1 * time.Minute) {
		expiration = 1 * time.Minute
	}
	result, err := c.db.SetNX(context.Background(), lockKey, value, expiration).Result()
	if err != nil {
		return false
	}
	return result
}

func (c *CacheService) ReleaseLock(key string, value string) bool {
	lockKey := fmt.Sprintf("lock:%s", key)
	existingValue := c.GetString(lockKey)
	if value == existingValue {
		return c.Delete(lockKey)
	}
	return false
}

func (c *CacheService) Delete(key string) bool {
	err := c.db.Del(context.Background(), key).Err()
	if err != nil {
		slog.Error("Error setting key in Cache", err)
	}
	return err != nil
}

func (c *CacheService) GetString(key string) string {
	value, err := c.db.Get(context.Background(), key).Result()
	switch {
	case err == redis.Nil:
		value = ""
	case err != nil:
		slog.Error("Error getting key from Cache", err)
	}
	return value
}

func (c *CacheService) SetString(key string, value string, expiration time.Duration) bool {
	err := c.db.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		slog.Error("Error setting key in Cache", err)
	}
	return err != nil
}

func (c *CacheService) GetBool(key string) bool {
	value, err := c.db.Get(context.Background(), key).Bool()
	switch {
	case err == redis.Nil:
		value = false
	case err != nil:
		slog.Error("Error getting key from Cache", err)
	}
	return value
}

func (c *CacheService) SetBool(key string, value bool, expiration time.Duration) bool {
	err := c.db.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		slog.Error("Error setting key in Cache", err)
	}
	return err != nil
}
