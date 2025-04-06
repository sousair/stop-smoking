package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sousair/gocore/pkg/event"
	"github.com/sousair/stop-smoking/internal/core"
)

func (u Usecase) UserMetadataSent(ctx context.Context, req *core.UserMetadataSentRequest) error {
	user, err := u.deps.UserRepo.FindOne(ctx, &core.User{
		Identifier: req.Identifier,
	})
	if err != nil {
		return err
	}

	userMetadata := &core.UserMetadata{
		WakeUpHour:             req.WakeUpHour,
		SleepHour:              req.SleepHour,
		SmokeToSleep:           req.SmokeToSleep,
		RolledOverSleepRoutine: req.SleepHour < req.WakeUpHour,
	}

	metadata, err := json.Marshal(userMetadata)
	if err != nil {
		return err
	}

	user.Metadata = metadata
	user.StartDate = time.Now()

	if _, err = u.deps.UserRepo.Update(ctx, user); err != nil {
		return err
	}

	return u.deps.EventEmitter.Emit(ctx,
		&event.Event{
			Type: core.StartProgramQueueEventType,
			Payload: &core.StartProgramScheduleRequest{
				Identifier:       user.Identifier,
				CigarettesPerDay: req.CigarettesPerDay,
			},
		},
		event.WithMaxRetries(3),
	)
}
