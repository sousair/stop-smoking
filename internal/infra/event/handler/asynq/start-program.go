package asynqhandler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sousair/stop-smoking/internal/core"
)

func (h AsynHandler) HandleStartProgramQueue(ctx context.Context, t *asynq.Task) error {
	var startReq *core.StartProgramScheduleRequest
	if err := json.Unmarshal(t.Payload(), &startReq); err != nil {
		return fmt.Errorf("[%s] failed to unmarshal payload: %v %w", t.Type(), err, asynq.SkipRetry)
	}

	return h.deps.DailyProgramUsecase.StartProgram(ctx, startReq)
}
