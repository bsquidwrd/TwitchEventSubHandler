package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type databaseService struct {
	db *pgxpool.Pool
}

func newDatabaseService() *databaseService {
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		panic(err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return &databaseService{
		db: db,
	}
}
