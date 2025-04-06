package usecase

import (
	"github.com/sousair/gocore/pkg/database/repository"
	"github.com/sousair/gocore/pkg/event"
	"github.com/sousair/stop-smoking/internal/core"
)

type Dependencies struct {
	DailyProgramRepo core.DailyProgramRepository
	UserRepo         repository.Repository[core.User]
	UserProgramRepo  repository.Repository[core.UserProgram]
	EventEmitter     event.EventEmitter
}

type Usecase struct {
	deps Dependencies
}

func New(deps Dependencies) *Usecase {
	return &Usecase{deps}
}
