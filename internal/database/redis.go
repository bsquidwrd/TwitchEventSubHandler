package database

import (
	"context"
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
		slog.Error("Error getting key from Redis", err)
	}
	return err != nil
}
