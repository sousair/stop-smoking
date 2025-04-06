package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/sousair/gocore/pkg/database"
	"github.com/sousair/gocore/pkg/database/entity"
	"github.com/sousair/gocore/pkg/event"
	"github.com/sousair/stop-smoking/internal/core"
)

func (u Usecase) CalculateNextProgram(ctx context.Context, userIdentifier int64) error {
	user, err := u.deps.UserRepo.FindOne(ctx, &core.User{
		Identifier: userIdentifier,
	})
	if err != nil {
		return err
	}

	userProgram, err := u.deps.UserProgramRepo.FindOne(ctx, &core.UserProgram{
		UserID:         user.ID,
		UserIdentifier: user.Identifier,
	})
	if err != nil {
		return err
	}

	if err := u.deps.UserProgramRepo.Delete(ctx, userProgram); err != nil {
		return err
	}

	program, err := u.deps.DailyProgramRepo.FindOne(ctx, &core.DailyProgram{
		BaseEntity: entity.BaseEntity{
			ID: userProgram.ProgramID,
		},
	})
	if err != nil {
		return err
	}

	if !userProgram.HasFailed {
		nextProgramID := userProgram.ProgramID - program.JumpOffset

		program, err := u.deps.DailyProgramRepo.FindOne(ctx, &core.DailyProgram{
			BaseEntity: entity.BaseEntity{
				ID: nextProgramID,
			},
		})
		if err != nil {
			if !errors.Is(err, database.ErrNotFound) {
				return err
			}

			// TODO: This is something that should alert admins
			return err
		}

		if _, err = u.deps.UserProgramRepo.Create(ctx, &core.UserProgram{
			UserID:         user.ID,
			UserIdentifier: user.Identifier,
			ProgramID:      program.ID,
		}); err != nil {
			return err
		}
	}

	var userMetadata *core.UserMetadata
	if err := json.Unmarshal(user.Metadata, &userMetadata); err != nil {
		return err
	}

	referenceTime := time.Now().Local()
	if !userMetadata.RolledOverSleepRoutine {
		referenceTime = time.Date(
			referenceTime.Year(), referenceTime.Month(), referenceTime.Day()+1,
			0, 0, 0, 0,
			referenceTime.Location(),
		)
	}

	schedule, err := program.CalculateUserSchedule(user, userProgram.ID, referenceTime)
	if err != nil {
		return err
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
