package currency

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/sladenkiy77/worker-service/internal/app/entity"
	"github.com/sladenkiy77/worker-service/internal/app/repository"
	"github.com/sladenkiy77/worker-service/internal/app/service"
)

// RatesProvider - источник "свежих" курсов (Fixer API).
// Интерфейс, а не конкретный клиент, чтобы сервис не зависел от реализации HTTP-клиента.
type RatesProvider interface {
	GetRates(ctx context.Context, base string) (map[string]float64, error)
}
