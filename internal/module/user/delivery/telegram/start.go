package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sousair/stop-smoking/internal/core"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
)

func (t *TelegramHandler) NewStartBotHandler(ctx context.Context) UpdateTelegramFunc {
	return func(update tgbotapi.Update) error {
		msg := update.Message

		identifier := msg.From.ID
		name := fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName)

		fmt.Printf("[StartHandler] Received [START] msg! [%d] %s\n", identifier, name)

		if err := t.deps.Usecase.CreateUser(ctx, &core.User{
			Identifier: identifier,
			Name:       name,
		}); err != nil {
			return err
		}

		formMsg := message.NewSimpleMessage(&message.Params{
			Content: `🚭 Welcome to Stop Smoking Bot! 🚭
This bot will help you to stop smoking by stricting your daily smoking limit and schedule.

Now we need to know some information about your daily smkoing habit.
Please answer the following form.
      `,
			Recipient: identifier,
			Bot:       t.deps.Bot,
		})

		if _, err := formMsg.Send(ctx, message.WithMessageButtons(message.KeyboardRow{
			{"Answer Form", StartFormKey},
		})); err != nil {
			return err
		}

		return nil
	}
}
