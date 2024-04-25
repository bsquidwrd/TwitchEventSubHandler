package database

import "github.com/jackc/pgx/v5/pgxpool"

type Service struct {
	Cache    *cacheService
	Database *pgxpool.Pool
	Twitch   *twitchService
}

func New() *Service {
	return &Service{
		Cache:    newCacheService(),
		Database: newDatabaseService(),
		Twitch:   newTwitchService(),
	}
}
