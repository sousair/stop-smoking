package asynqhandler

import (
	"github.com/sousair/gocore/pkg/event"
	dailyProgramUsecase "github.com/sousair/stop-smoking/internal/module/dailyprogram/usecase"
)

type Dependencies struct {
	DailyProgramUsecase *dailyProgramUsecase.Usecase
	EventEmitter        event.EventEmitter
}

type AsynHandler struct {
	deps Dependencies
}

func New(deps Dependencies) AsynHandler {
	return AsynHandler{deps: deps}
}
