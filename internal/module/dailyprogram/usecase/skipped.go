package usecase

import (
	"context"

	"github.com/sousair/stop-smoking/internal/core"
)

func (u Usecase) Skipped(ctx context.Context, identifier int64) error {
	userProgram, err := u.deps.UserProgramRepo.FindOne(ctx, &core.UserProgram{
		UserIdentifier: identifier,
	})
	if err != nil {
		return err
	}

	userProgram.Skipped++

	if _, err = u.deps.UserProgramRepo.Update(ctx, userProgram); err != nil {
		return err
	}

	return nil
}
