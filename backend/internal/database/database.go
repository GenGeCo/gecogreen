package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// DB holds database connections
type DB struct {
	Pool  *pgxpool.Pool
	Redis *redis.Client
}

// New creates a new database connection
func New(databaseURL, redisURL string) (*DB, error) {
	// PostgreSQL connection pool
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Connection pool settings
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test PostgreSQL connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	// Redis connection
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to parse redis URL: %w", err)
	}

	redisClient := redis.NewClient(opt)

	// Test Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping redis: %w", err)
	}

	return &DB{
		Pool:  pool,
		Redis: redisClient,
	}, nil
}

// Close closes all database connections
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
	if db.Redis != nil {
		db.Redis.Close()
	}
}

// Health checks database health
func (db *DB) Health(ctx context.Context) map[string]string {
	health := make(map[string]string)

	// Check PostgreSQL
	if err := db.Pool.Ping(ctx); err != nil {
		health["postgres"] = fmt.Sprintf("error: %v", err)
	} else {
		health["postgres"] = "ok"
	}

	// Check Redis
	if _, err := db.Redis.Ping(ctx).Result(); err != nil {
		health["redis"] = fmt.Sprintf("error: %v", err)
	} else {
		health["redis"] = "ok"
	}

	return health
}
