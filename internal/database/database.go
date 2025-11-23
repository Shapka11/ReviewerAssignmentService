package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"ReviewerAssignmentService/internal/config"
)

func NewDBPool(cfg *config.Config) (*pgxpool.Pool, error) {
	cfgString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	poolConfig, err := pgxpool.ParseConfig(cfgString)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = 10

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database PostgresSQL")
	return pool, nil
}
