package telegram

import (
	"github.com/sousair/gocore/pkg/cache"
	"github.com/sousair/stop-smoking/internal/module/dailyprogram/usecase"
	tbot "github.com/thunderjr/go-telegram/pkg/bot"
)

const CanSmokeKeyPattern = "%d:can-smoke"

type Dependencies struct {
	Usecase *usecase.Usecase
	Bot     *tbot.TelegramBot
	Cache   cache.Cache
}

type TelegramHandler struct {
	deps Dependencies
}

func New(deps Dependencies) *TelegramHandler {
	return &TelegramHandler{deps}
}
