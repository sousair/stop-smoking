package asynqhandler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sousair/stop-smoking/internal/core"
)

func (h AsynHandler) HandleCalculateNextProgramQueue(ctx context.Context, t *asynq.Task) error {
	var schedule *core.SmokingSchedule
	if err := json.Unmarshal(t.Payload(), &schedule); err != nil {
		return fmt.Errorf("[%s] failed to unmarshal payload: %v %w", t.Type(), err, asynq.SkipRetry)
	}

	return h.deps.DailyProgramUsecase.CalculateNextProgram(ctx, schedule.Identifier)
}
