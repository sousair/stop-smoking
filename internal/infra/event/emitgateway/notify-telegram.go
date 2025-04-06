package emitgateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sousair/gocore/pkg/event"
)

var (
	NotifyTelegramHandlerNotFound = errors.New("NotifyTelegram handler not found")
)

func (e *EventEmit) NewNotifyTelegramEventEmitter() event.EmitHandler {
	return func(ctx context.Context, event *event.Event, opt *event.EventOption) error {
		handler, ok := e.TelegramHandlers[event.Type]
		if !ok {
			return NotifyTelegramHandlerNotFound
		}

		data, err := json.Marshal(event.Payload)
		if err != nil {
			return err
		}

		fn := func() error {
			return handler(ctx, data)
		}

		if opt.MaxRetries > 0 {
			fn = func() error {
				sleep := 1 * time.Second
				for i := 0; i < opt.MaxRetries; i++ {
					if i > 0 {
						fmt.Printf("[%s] Retrying... \n", event.Type)
						time.Sleep(sleep)
						sleep *= 2
					}

					err := handler(ctx, data)
					if err == nil {
						return nil
					}

					fmt.Printf("[%s] Error: %v \n", event.Type, err)
				}

				return fmt.Errorf("[%s] failed to emit event: %w", event.Type, err)
			}
		}

		if !opt.ProcessAt.IsZero() {
			go func() {
				time.Sleep(time.Until(opt.ProcessAt))
				if err := fn(); err != nil {
					fmt.Printf("[%s] Error: %v \n", event.Type, err)
				}
			}()

			return nil
		}

		if opt.Delay > 0 {
			go func() {
				time.Sleep(opt.Delay)
				if err := fn(); err != nil {
					fmt.Printf("[%s] Error: %v", event.Type, err)
				}
			}()

			return nil
		}

		return fn()
	}
}
