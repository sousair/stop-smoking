package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/sousair/stop-smoking/internal/core"
	"github.com/thunderjr/go-telegram/pkg/bot/update"
)

func (t TelegramHandler) NewStartProgramFormFilled(ctx context.Context) UpdateFormTelegramFunc[StartProgramForm] {
	return func(ctx context.Context, form *update.FormAnswerData[StartProgramForm]) error {
		lastAnswerMessage := form.LastAnswerUpdate.Message
		identifier := lastAnswerMessage.Chat.ID

		fmt.Printf(
			"[StartProgramFormHandler] Received [START_PROGRAM_FORM] msg! [%d] %s \n",
			identifier,
			lastAnswerMessage.From.FirstName,
		)

		formData := form.Data

		if err := t.deps.Usecase.StartProgram(ctx, &core.StartProgramRequest{
			Identifier:       identifier,
			CigarettesPerDay: formData.CigarettesPerDay,
			WakeUpHour:       formData.WakeUpHour,
			SleepHour:        formData.SleepHour,
			SmokeToSleep:     strings.ToLower(formData.SmokeToSleep) == "yes",
		}); err != nil {
			fmt.Printf("error aqui oh %v \n", err)
			return err
		}

		return nil
	}
}
