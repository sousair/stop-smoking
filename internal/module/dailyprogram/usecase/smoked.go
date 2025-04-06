package usecase

import (
	"context"

	"github.com/sousair/stop-smoking/internal/core"
)

func (u Usecase) Smoked(ctx context.Context, req *core.SmokedRequest) error {
	userProgram, err := u.deps.UserProgramRepo.FindOne(ctx, &core.UserProgram{
		UserIdentifier: req.Identifier,
	})
	if err != nil {
		return err
	}

	userProgram.Smoked++

	if req.Failed {
		userProgram.Failed++
		if !userProgram.HasFailed && userProgram.Skipped < userProgram.Failed {
			userProgram.HasFailed = true
		}
	}

	if _, err = u.deps.UserProgramRepo.Update(ctx, userProgram); err != nil {
		return err
	}
	return nil
}
