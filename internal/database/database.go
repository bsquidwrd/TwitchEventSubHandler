package database

type Service struct {
	Redis  *redisService
	Twitch *twitchService
}

func New() *Service {
	return &Service{
		Redis:  newRedisService(),
		Twitch: newTwitchService(),
	}
}
