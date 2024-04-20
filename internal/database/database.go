package database

import (
	"os"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	Redis *redisService
}

func New() *Services {
	redisUrl := os.Getenv("REDIS_CONNECTIONSTRING")
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	return &Services{
		Redis: &redisService{db: redis.NewClient(opt)},
	}
}
