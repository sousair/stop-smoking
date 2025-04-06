package emitgateway

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
	"github.com/sousair/gocore/pkg/event"
)

func (e *EventEmit) NewQueueAsynqEventEmitter() event.EmitHandler {
	return func(ctx context.Context, event *event.Event, opt *event.EventOption) error {
		data, err := json.Marshal(event.Payload)
		if err != nil {
			return err
		}

		task := asynq.NewTask(
			event.Type.String(),
			data,
			getAsynqOptions(opt)...,
		)

		info, err := e.deps.AsynqClient.EnqueueContext(ctx, task)
		if err != nil {
			return err
		}

		log.Printf("[AsynqEmitter] task enqueued: id=%s queue=%s", info.ID, info.Queue)
		return nil
	}
}

func getAsynqOptions(opt *event.EventOption) []asynq.Option {
	options := []asynq.Option{}

	if opt.MaxRetries > 0 {
		options = append(options, asynq.MaxRetry(opt.MaxRetries))
	}
	if opt.Delay > 0 {
		options = append(options, asynq.ProcessIn(opt.Delay))
	}
	if !opt.ProcessAt.IsZero() {
		options = append(options, asynq.ProcessAt(opt.ProcessAt))
	}
	return options
}
