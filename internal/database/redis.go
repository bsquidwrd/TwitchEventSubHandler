package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisService struct {
	db *redis.Client
}

func newRedisService() *redisService {
	redisUrl := os.Getenv("REDIS_CONNECTIONSTRING")
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}
	slog.Info("Redis connected successfully")

	client := redis.NewClient(opt)
	_, err = client.Get(context.Background(), "thiswillmostlikelyneverexistandthatsokay").Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	return &redisService{
		db: client,
	}
}

func (r *redisService) TakeLock(key string, value string, expiration time.Duration) bool {
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

func (r *redisService) ReleaseLock(key string, value string) bool {
	lockKey := fmt.Sprintf("lock:%s", key)
	existingValue := r.GetString(lockKey)
	if value == existingValue {
		return r.Delete(lockKey)
	}
	return false
}

func (r *redisService) Delete(key string) bool {
	err := r.db.Del(context.Background(), key).Err()
	if err != nil {
		slog.Error("Error settings key in Redis", err)
	}
	return err != nil
}

func (r *redisService) GetString(key string) string {
	value, err := r.db.Get(context.Background(), key).Result()
	switch {
	case err == redis.Nil:
		value = ""
	case err != nil:
		slog.Error("Error getting key from Redis", err)
	}
	return value
}

func (r *redisService) SetString(key string, value string, expiration time.Duration) bool {
	err := r.db.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		slog.Error("Error settings key in Redis", err)
	}
	return err != nil
}

func (r *redisService) GetBool(key string) bool {
	value, err := r.db.Get(context.Background(), key).Bool()
	switch {
	case err == redis.Nil:
		value = false
	case err != nil:
		slog.Error("Error getting key from Redis", err)
	}
	return value
}

func (r *redisService) SetBool(key string, value bool, expiration time.Duration) bool {
	err := r.db.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		slog.Error("Error settings key in Redis", err)
	}
	return err != nil
}
