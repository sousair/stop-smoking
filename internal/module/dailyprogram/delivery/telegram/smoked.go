package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sousair/stop-smoking/internal/core"
)

func (h TelegramHandler) NewSmokedHandler(
	ctx context.Context,
) func(update tgbotapi.Update) error {
	return func(update tgbotapi.Update) error {
		var identifier int64

		if update.Message != nil {
			identifier = update.Message.Chat.ID
		} else {
			identifier = update.CallbackQuery.Message.Chat.ID
		}

		cacheKey := fmt.Sprintf(CanSmokeKeyPattern, identifier)

		failed := true
		if val, err := h.deps.Cache.Get(ctx, cacheKey); err == nil {
			if val == "true" {
				failed = false
			}
		}

		if err := h.deps.Usecase.Smoked(ctx, &core.SmokedRequest{
			Identifier: identifier,
			Failed:     failed,
		}); err != nil {
			return err
		}

		return h.deps.Cache.Delete(ctx, cacheKey)
	}
}
