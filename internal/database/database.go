package database

import "github.com/jackc/pgx/v5/pgxpool"

type Service struct {
	Cache    *cacheService
	Database *pgxpool.Pool
	Twitch   *twitchService
	Queue    *queueService
}

func New() *Service {
	return &Service{
		Cache:    newCacheService(),
		Database: newDatabaseService(),
		Twitch:   newTwitchService(),
		Queue:    newQueueService(),
	}
}

func (s *Service) Cleanup() {
	defer s.Cache.cleanup()
	defer s.Database.Close()
	defer s.Queue.cleanup()
}
