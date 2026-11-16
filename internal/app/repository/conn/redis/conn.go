package rcredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/sladenkiy77/worker-service/internal/app/config/section"
)

type Client struct {
	*redis.Client
}

func NewClient(ctx context.Context, cfg section.RepositoryRedis) (*Client, error) {
	log.Info().Str("address", cfg.Address).Int("db_num", cfg.DB).Msg("started connection to redis db")
	rdb := redis.NewClient(&redis.Options{Addr: cfg.Address, Password: cfg.Password, DB: cfg.DB})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("failed to start redis tracing: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Info().Str("address", cfg.Address).Int("db_num", cfg.DB).Msg("successfully connected to redis db")

	return &Client{Client: rdb}, nil
}
