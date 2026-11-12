package eorder

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/sladenkiy77/worker-service/internal/app/entity"
	ehandler "github.com/sladenkiy77/worker-service/internal/app/handler/event"
	"github.com/sladenkiy77/worker-service/internal/pkg/broker"
)

type handler struct{}

func NewHandler() ehandler.OrderCreated {
	return &handler{}
}

func (h *handler) CallbackOrderCreated(ctx context.Context, ev *entity.EventOrderCreated, _ []broker.Header,
) error {
	log.Info().Ctx(ctx).EmbedObject(ev).Msg("received order created event")

	if ev.OrderGUID == "" {
		return broker.NotCriticalError(
			fmt.Errorf("order_guid is empty"),
		)
	}

	return nil
}
