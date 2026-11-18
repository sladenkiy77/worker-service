package currency

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/sladenkiy77/worker-service/internal/app/entity"
	"github.com/sladenkiy77/worker-service/internal/app/repository"
)

const rateKeyPrefix = "rates:"

type repoRedis struct {
	client   *redis.Client
	cacheTTL time.Duration
}

func NewRepoFromRedis(client *redis.Client, cacheTTL time.Duration) repository.CurrencyRate {
	return &repoRedis{client: client, cacheTTL: cacheTTL}
}

func (r *repoRedis) buildKey(from, to string) string {
	return fmt.Sprintf("%s%s:%s", rateKeyPrefix, from, to)
}

func (r *repoRedis) GetRate(ctx context.Context, from, to string) (float64, error) {
	rate, err := r.client.Get(ctx, r.buildKey(from, to)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, entity.ErrRateNotFound
		}

		return 0, err
	}

	res, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return 0, fmt.Errorf("parse cached rate %q: %w", rate, err)
	}

	return res, nil
}

func (r *repoRedis) SetRate(ctx context.Context, from, to string, rate float64) error {
	return r.client.Set(
		ctx,
		r.buildKey(from, to),
		strconv.FormatFloat(rate, 'f', -1, 64),
		r.cacheTTL).
		Err()
}

func (r *repoRedis) SetRates(ctx context.Context, from string, rates map[string]float64) error {
	pipe := r.client.Pipeline()

	for to, rate := range rates {
		pipe.Set(
			ctx,
			r.buildKey(from, to),
			strconv.FormatFloat(rate, 'f', -1, 64),
			r.cacheTTL)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipeline rates set exec: %w", err)
	}

	log.Info().Ctx(ctx).Str("from", from).Int("count", len(rates)).Msg("Курсы валют сохранены в кэш")
	return nil
}
