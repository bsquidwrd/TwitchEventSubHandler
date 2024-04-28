package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DiscordNotifierService struct {
	Cache    *cacheService
	Database *pgxpool.Pool
	Queue    *queueService
}

func NewDiscordNotifierService() *DiscordNotifierService {
	return &DiscordNotifierService{
		Cache:    newCacheService(),
		Database: newDatabaseService(),
		Queue:    newQueueService(),
	}
}

func (s *DiscordNotifierService) Cleanup() {
	defer s.Cache.cleanup()
	defer s.Database.Close()
	defer s.Queue.cleanup()
}

func (s *DiscordNotifierService) HealthCheck() error {
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
	return nil
}
