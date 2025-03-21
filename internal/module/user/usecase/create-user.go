package usecase

import (
	"context"
	"errors"

	"github.com/sousair/gocore/pkg/database"
	"github.com/sousair/stop-smoking/internal/core"
)

var (
	UserAlreadyExistsErr = errors.New("user already exists")
)

func (u Usecase) CreateUser(ctx context.Context, req *core.User) error {
	user, err := u.deps.UserRepo.FindOne(ctx, &core.User{
		Idendifier: req.Idendifier,
	})
	if err != nil {
		if err != database.ErrNotFound {
			return err
		}
	}

	if user != nil {
		return UserAlreadyExistsErr
	}

	_, err = u.deps.UserRepo.Create(ctx, req)

	return err
}
