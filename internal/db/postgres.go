package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultPingAttempts = 5
	defaultPingTimeout  = 3 * time.Second
	defaultPingDelay    = 2 * time.Second
)

func NewPostgresPool(ctx context.Context, connURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connURL)
	if err != nil {
		return nil, fmt.Errorf("parse connection url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	err = pingWithRetry(ctx, pool)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("connection to postgres: %w", err)
	}

	log.Println("LOG - database connection pool established")
	return pool, nil
}

func pingWithRetry(ctx context.Context, pool *pgxpool.Pool) error {
	var pingErr error
	for attempt := 1; attempt <= defaultPingAttempts; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		pingCtx, cancel := context.WithTimeout(ctx, defaultPingTimeout)
		pingErr = pool.Ping(pingCtx)
		cancel()

		if pingErr == nil {
			log.Printf("LOG - connected to postgres on attempt %d", attempt)
			return nil
		}

		log.Printf("LOG - postgres not ready (attempt %d/%d): %v", attempt, defaultPingAttempts, pingErr)

		select {
		case <-time.After(defaultPingDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return pingErr
}
