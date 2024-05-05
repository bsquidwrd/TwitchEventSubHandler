package service

import (
	"context"
	"log/slog"

	"github.com/bsquidwrd/TwitchEventSubHandler/shared/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReceiverService struct {
	Cache    *database.CacheService
	Database *pgxpool.Pool
	Twitch   *twitchService
	Queue    *queueService
}

func NewReceiverService() *ReceiverService {
	return &ReceiverService{
		Cache:    database.NewCacheService(),
		Database: database.NewDatabaseService(),
		Twitch:   newTwitchService(),
		Queue:    newQueueService(),
	}
}

func (s *ReceiverService) Cleanup() {
	defer s.Cache.Cleanup()
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
