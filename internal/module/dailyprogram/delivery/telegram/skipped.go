package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h TelegramHandler) NewSkippedHandler(
	ctx context.Context,
) func(update tgbotapi.Update) error {
	return func(update tgbotapi.Update) error {
		identifier := update.CallbackQuery.Message.Chat.ID

		cacheKey := fmt.Sprintf(CanSmokeKeyPattern, identifier)

		val, err := h.deps.Cache.Get(ctx, cacheKey)
		if err != nil || val == "false" {
			return nil
		}

		if err := h.deps.Usecase.Skipped(ctx, identifier); err != nil {
			return err
		}

		return h.deps.Cache.Delete(ctx, cacheKey)
	}
}
