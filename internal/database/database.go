package database

type Services struct {
	Redis  *redisService
	Twitch *twitchService
}

func New() *Services {
	return &Services{
		Redis:  newRedisService(),
		Twitch: newTwitchService(),
	}
}
