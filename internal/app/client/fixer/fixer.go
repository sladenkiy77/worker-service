package fixer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/sladenkiy77/worker-service/internal/app/config/section"
	"github.com/sladenkiy77/worker-service/internal/app/entity"
)

const requestTimeout = 10 * time.Second

// Client - HTTP-клиент Fixer API. Транспорт обернут otelhttp, пожтому исходящие запросы трассируются (no-op при выключенном OpenTelemetry)
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

type fixerResponse struct {
	Success bool               `json:"success"`
	Base    string             `json:"base"`
	Date    string             `json:"date"`
	Rates   map[string]float64 `json:"rates"`
	Error   *fixerError        `json:"error,omitempty"`
}

type fixerError struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}
