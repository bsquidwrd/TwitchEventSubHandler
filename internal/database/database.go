package database

type Services struct {
	Redis *redisService
}

func New() *Services {
	return &Services{
		Redis: newRedisService(),
	}
}
