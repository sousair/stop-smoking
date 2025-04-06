package emitgateway

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/sousair/gocore/pkg/event"
)

type TelegramHandler func(ctx context.Context, payload []byte) error

type Dependencies struct {
	AsynqClient *asynq.Client
}

type EventEmit struct {
	deps             Dependencies
	TelegramHandlers map[event.EventType]TelegramHandler
}

func New(deps Dependencies) *EventEmit {
	return &EventEmit{
		deps:             deps,
		TelegramHandlers: make(map[event.EventType]TelegramHandler),
	}
}

func (e *EventEmit) AddTelegramHandlers(handlers map[event.EventType]TelegramHandler) {
	for key, handler := range handlers {
		e.TelegramHandlers[key] = handler
	}
}
