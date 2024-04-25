package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type cacheService struct {
	db *redis.Client
}

func newCacheService() *cacheService {
	cacheUrl := os.Getenv("CACHE_CONNECTIONSTRING")
	opt, err := redis.ParseURL(cacheUrl)
	if err != nil {
		panic(err)
	}
	slog.Info("Cache connected successfully")

	client := redis.NewClient(opt)
	_, err = client.Get(context.Background(), "thiswillmostlikelyneverexistandthatsokay").Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	return &cacheService{
		db: client,
	}
}

func (r *cacheService) TakeLock(key string, value string, expiration time.Duration) bool {
	lockKey := fmt.Sprintf("lock:%s", key)
	// Have a reasonable default a lock can be taken for
	if expiration > (1 * time.Minute) {
		expiration = 1 * time.Minute
	}
	result, err := r.db.SetNX(context.Background(), lockKey, value, expiration).Result()
	if err != nil {
		return false
	}
	return result
}

func (r *cacheService) ReleaseLock(key string, value string) bool {
	lockKey := fmt.Sprintf("lock:%s", key)
	existingValue := r.GetString(lockKey)
	if value == existingValue {
		return r.Delete(lockKey)
	}
	return false
}

func (r *cacheService) Delete(key string) bool {
	err := r.db.Del(context.Background(), key).Err()
	if err != nil {
		slog.Error("Error setting key in Cache", err)
	}
	return err != nil
}

func (r *cacheService) GetString(key string) string {
	value, err := r.db.Get(context.Background(), key).Result()
	switch {
	case err == redis.Nil:
		value = ""
	case err != nil:
		slog.Error("Error getting key from Cache", err)
	}
	return value
}

func (r *cacheService) SetString(key string, value string, expiration time.Duration) bool {
	err := r.db.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		slog.Error("Error setting key in Cache", err)
	}
	return err != nil
}

func (r *cacheService) GetBool(key string) bool {
	value, err := r.db.Get(context.Background(), key).Bool()
	switch {
	case err == redis.Nil:
		value = false
	case err != nil:
		slog.Error("Error getting key from Cache", err)
	}
	return value
}

func (r *cacheService) SetBool(key string, value bool, expiration time.Duration) bool {
	err := r.db.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		slog.Error("Error setting key in Cache", err)
	}
	return err != nil
}
