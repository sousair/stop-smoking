package telegram

import "github.com/thunderjr/go-telegram/pkg/bot/update"

type StartProgramForm struct {
	CigarettesPerDay int
	WakeUpHour       int
	SleepHour        int
	SmokeToSleep     string
}

var _ update.PromptProvider = (*StartProgramForm)(nil)

func (s *StartProgramForm) FieldPrompts() ([]update.FormFieldPrompt, error) {
	return []update.FormFieldPrompt{
		{
			Name:   "CigarettesPerDay",
			Prompt: "How many cigarettes do you smoke per day?",
			Order:  1,
		},
		{
			Name:   "WakeUpHour",
			Prompt: "What time do you usually wake up? [0-23]",
			Order:  2,
		},
		{
			Name:   "SleepHour",
			Prompt: "What time do you usually go to sleep? [0-23]",
			Order:  3,
		},
		{
			Name:   "SmokeToSleep",
			Prompt: "Do you NEED to smoke before going to sleep? [yes/no]",
			Order:  4,
		},
	}, nil
}
