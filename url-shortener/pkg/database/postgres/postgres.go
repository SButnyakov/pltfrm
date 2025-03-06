package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"url-shortener/pkg/config"
)

func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	connectUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name)

	pool, err := pgxpool.New(context.Background(), connectUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return pool, nil
}
