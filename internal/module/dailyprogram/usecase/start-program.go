package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sousair/gocore/pkg/event"
	"github.com/sousair/stop-smoking/internal/core"
)

func (u Usecase) StartProgram(ctx context.Context, req *core.StartProgramScheduleRequest) error {
	user, err := u.deps.UserRepo.FindOne(ctx, &core.User{
		Identifier: req.Identifier,
	})
	if err != nil {
		return err
	}

	var userMetadata *core.UserMetadata
	if err := json.Unmarshal(user.Metadata, &userMetadata); err != nil {
		return err
	}

	program, err := u.deps.DailyProgramRepo.FindOne(ctx, &core.DailyProgram{
		CigarretesPerDay: req.CigarettesPerDay,
	})
	if err != nil {
		return err
	}

	userProgram, err := u.deps.UserProgramRepo.Create(ctx, &core.UserProgram{
		UserID:         user.ID,
		UserIdentifier: user.Identifier,
		ProgramID:      program.ID,
	})
	if err != nil {
		return err
	}

	var schedule *core.SmokingSchedule
	referenceTime := time.Now().Local()
	schedule, err = program.CalculateUserSchedule(user, userProgram.ID, referenceTime)
	// WARN: This is ugly.
	if err != nil {
		schedule, err = program.CalculateUserSchedule(user, userProgram.ID, time.Date(
			referenceTime.Year(), referenceTime.Month(), referenceTime.Day()+1,
			0, 0, 0, 0,
			referenceTime.Location(),
		))
		if err != nil {
			return err
		}
	}

	return u.deps.EventEmitter.Emit(ctx,
		&event.Event{
			Type:    core.ScheduleSmokingTimeQueueEventType,
			Payload: schedule,
		},
		event.WithProcessAt(schedule.NextSmoke),
		event.WithMaxRetries(3),
	)
}
