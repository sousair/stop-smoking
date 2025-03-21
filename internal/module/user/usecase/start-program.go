package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sousair/stop-smoking/internal/core"
)

func (u Usecase) StartProgram(ctx context.Context, req *core.StartProgramRequest) error {
	user, err := u.deps.UserRepo.FindOne(ctx, &core.User{
		Idendifier: req.Identifier,
	})
	if err != nil {
		return err
	}

	program, err := u.deps.DailyProgramRepo.FindProgramByPerDay(ctx, req.CigarettesPerDay)
	if err != nil {
		return err
	}

	userMetadata := &core.UserMetadata{
		WakeUpHour:   req.WakeUpHour,
		SleepHour:    req.SleepHour,
		SmokeToSleep: req.SmokeToSleep,
	}

	metadata, err := json.Marshal(userMetadata)
	if err != nil {
		return err
	}

	user.ProgramID = program.ID
	user.Metadata = metadata
	user.StartDate = time.Now()

	if _, err = u.deps.UserRepo.Update(ctx, user); err != nil {
		return err
	}

	// TODO: Call to schedule program with `go`

	return nil
}
