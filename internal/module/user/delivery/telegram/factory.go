package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sousair/stop-smoking/internal/module/user/usecase"
	tbot "github.com/thunderjr/go-telegram/pkg/bot"
	"github.com/thunderjr/go-telegram/pkg/bot/update"
)

type Dependencies struct {
	Usecase *usecase.Usecase
	Bot     *tbot.TelegramBot
}

type TelegramHandler struct {
	deps Dependencies
}

type UpdateTelegramFunc func(update tgbotapi.Update) error
type UpdateFormTelegramFunc[T any] func(ctx context.Context, form *update.FormAnswerData[T]) error

func New(deps Dependencies) *TelegramHandler {
	return &TelegramHandler{deps}
}
