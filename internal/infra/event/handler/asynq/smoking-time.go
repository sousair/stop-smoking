package asynqhandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sousair/gocore/pkg/event"
	"github.com/sousair/stop-smoking/internal/core"
)

func (h AsynHandler) HandleSmokingTimeQueue(ctx context.Context, t *asynq.Task) error {
	var schedule *core.SmokingSchedule
	if err := json.Unmarshal(t.Payload(), &schedule); err != nil {
		return fmt.Errorf("[%s] failed to unmarshal payload: %v %w", t.Type(), err, asynq.SkipRetry)
	}

	if err := h.deps.EventEmitter.Emit(ctx,
		&event.Event{
			Type:    core.NotifySmokingTimeEventType,
			Payload: schedule,
		},
		event.WithMaxRetries(3),
	); err != nil {
		return err
	}

	if err := schedule.Next(); err != nil {
		if !errors.Is(err, core.NoSmokingTimeAvailableErr) {
			return err
		}

		return h.deps.EventEmitter.Emit(ctx,
			&event.Event{
				Type:    core.CalculateNextProgramQueueEventType,
				Payload: schedule,
			},
			event.WithDelay(time.Minute*10),
			event.WithMaxRetries(3),
		)
	}

	return h.deps.EventEmitter.Emit(ctx,
		&event.Event{
			Type:    core.ScheduleSmokingTimeQueueEventType,
			Payload: schedule,
		},
		event.WithProcessAt(schedule.NextSmoke),
		event.WithMaxRetries(3),
	)
}
