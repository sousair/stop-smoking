package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sousair/gocore/pkg/database"
	"github.com/sousair/gocore/pkg/database/repository"
	"github.com/sousair/stop-smoking/internal/core"
	dailyprogramRepo "github.com/sousair/stop-smoking/internal/module/dailyprogram/repository"
	userTelegram "github.com/sousair/stop-smoking/internal/module/user/delivery/telegram"
	userUsecase "github.com/sousair/stop-smoking/internal/module/user/usecase"
	tbot "github.com/thunderjr/go-telegram/pkg/bot"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
	"github.com/thunderjr/go-telegram/pkg/bot/update"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()

	tbot.SetAppName("stop-smoking")

	ctx = tbot.WithReplyActionRepo(ctx, tbot.NewRepository[message.ReplyAction]())
	ctx = tbot.WithFormAnswerRepo(ctx, tbot.NewRepository[update.FormAnswer]())

	db, err := database.NewSQLite()
	handleErr(err)

	userRepo, err := repository.NewRepository[core.User](db)
	handleErr(err)

	dailyProgramRepo, err := dailyprogramRepo.NewDailyProgramRepository(db)
	handleErr(err)

	userUsecase := userUsecase.New(userUsecase.Dependencies{
		UserRepo:         userRepo,
		DailyProgramRepo: dailyProgramRepo,
	})

	bot, err := tbot.New(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
		tbot.WithUpdateHandlers([]update.Handler{}),
	)
	handleErr(err)

	telegramHandler := userTelegram.New(userTelegram.Dependencies{
		Usecase: userUsecase,
		Bot:     bot,
	})

	bot.AddHandlers(
		append(
			update.NewFormHandlers(ctx, &update.NewFormHandlerParams[userTelegram.StartProgramForm]{
				Type:     update.HandlerTypeKeyboardCallback,
				Key:      userTelegram.StartFormKey,
				Form:     &userTelegram.StartProgramForm{},
				OnSubmit: telegramHandler.NewStartProgramFormFilled(ctx),
				Bot:      bot,
			}),
			update.NewMessageUpdate(userTelegram.StartKey, telegramHandler.NewStartBotHandler(ctx)),
		)...,
	)

	errChan := make(chan error)

	go func() {
		for err := range errChan {
			fmt.Printf("Received error: %v\n", err)
		}
	}()

	bot.Updates(ctx, errChan)
}
