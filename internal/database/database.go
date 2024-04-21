package database

import "sync"

type Services struct {
	Redis    *redisService
	AuthLock *sync.Mutex
}

func New() *Services {
	return &Services{
		Redis:    newRedisService(),
		AuthLock: &sync.Mutex{},
	}
}
