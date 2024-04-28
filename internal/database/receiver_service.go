package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReceiverService struct {
	Cache    *cacheService
	Database *pgxpool.Pool
	Twitch   *twitchService
	Queue    *queueService
}

func NewReceiverService() *ReceiverService {
	return &ReceiverService{
		Cache:    newCacheService(),
		Database: newDatabaseService(),
		Twitch:   newTwitchService(),
		Queue:    newQueueService(),
	}
}

func (s *ReceiverService) Cleanup() {
	defer s.Cache.cleanup()
	defer s.Database.Close()
	defer s.Queue.cleanup()
}

func (s *ReceiverService) HealthCheck() error {
	var err error

	err = s.Cache.Ping()
	if err != nil {
		slog.Error("Cache failed healthcheck", "error", err)
		return err
	}

	err = s.Database.Ping(context.Background())
	if err != nil {
		slog.Error("Database failed healthcheck", "error", err)
		return err
	}

	err = s.Queue.Ping()
	if err != nil {
		slog.Error("Queue failed healthcheck", "error", err)
		return err
	}

	err = s.Twitch.Ping()
	if err != nil {
		slog.Error("Twitch failed healthcheck", "error", err)
		return err
	}
	return nil
}
