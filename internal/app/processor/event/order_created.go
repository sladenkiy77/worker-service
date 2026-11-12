package eprocessor

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/sladenkiy77/worker-service/internal/app/entity"
	ehandler "github.com/sladenkiy77/worker-service/internal/app/handler/event"
	"github.com/sladenkiy77/worker-service/internal/app/processor"
	"github.com/sladenkiy77/worker-service/internal/pkg/broker"
)

type orderCreatedProc struct {
	h   ehandler.OrderCreated
	bus broker.Bus[entity.EventOrderCreated]
}

func NewOrderCreatedEventsCatcher(
	h ehandler.OrderCreated,
	bus broker.Bus[entity.EventOrderCreated],
) processor.Processor {
	return &orderCreatedProc{h, bus}
}

func (p *orderCreatedProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	if err := p.bus.Subscribe(ctx, wg, p.h.CallbackOrderCreated); err != nil {
		log.Fatal().Err(err).Str("topic", p.bus.QueueName()).Msg("failed to subscribe to topic")
	}

	log.Info().Str("topic", p.bus.QueueName()).Msg("successfully subscribed to topic")
}
