package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sousair/gocore/pkg/cache"
	"github.com/sousair/gocore/pkg/database"
	"github.com/sousair/gocore/pkg/database/repository"
	"github.com/sousair/gocore/pkg/event"
	"github.com/sousair/stop-smoking/internal/core"
	"github.com/sousair/stop-smoking/internal/infra/event/emitgateway"
	asynqhandler "github.com/sousair/stop-smoking/internal/infra/event/handler/asynq"
	dailyprogramTelegram "github.com/sousair/stop-smoking/internal/module/dailyprogram/delivery/telegram"
	dailyprogramRepo "github.com/sousair/stop-smoking/internal/module/dailyprogram/repository"
	dailyprogramUsecase "github.com/sousair/stop-smoking/internal/module/dailyprogram/usecase"
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

var err error
var redisUrl = fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

func main() {
	ctx := context.Background()

	tbot.SetAppName("stop-smoking")

	ctx = tbot.WithReplyActionRepo(ctx, tbot.NewRepository[message.ReplyAction]())
	ctx = tbot.WithFormAnswerRepo(ctx, tbot.NewRepository[update.FormAnswer]())

	db, err := database.NewSQLite()
	handleErr(err)

	cache, err := cache.New()
	handleErr(err)

	userRepo, err := repository.NewRepository[core.User](db)
	handleErr(err)

	userProgramRepo, err := repository.NewRepository[core.UserProgram](db)
	handleErr(err)

	dailyProgramRepo, err := dailyprogramRepo.NewDailyProgramRepository(db)
	handleErr(err)

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisUrl,
	})

	emitGateway := emitgateway.New(emitgateway.Dependencies{
		AsynqClient: asynqClient,
	})

	eventEmitter := event.New(
		event.WithPrefixEmitHandler("notify",
			emitGateway.NewNotifyTelegramEventEmitter(),
		),
		event.WithPrefixEmitHandler("queue",
			emitGateway.NewQueueAsynqEventEmitter(),
		),
	)

	userUsecase := userUsecase.New(userUsecase.Dependencies{
		UserRepo:         userRepo,
		DailyProgramRepo: dailyProgramRepo,
		EventEmitter:     eventEmitter,
	})

	dailyProgramUsecase := dailyprogramUsecase.New(dailyprogramUsecase.Dependencies{
		DailyProgramRepo: dailyProgramRepo,
		UserRepo:         userRepo,
		UserProgramRepo:  userProgramRepo,
		EventEmitter:     eventEmitter,
	})

	bot, err := tbot.New(os.Getenv("TELEGRAM_BOT_TOKEN"))
	handleErr(err)

	userTelegramD := userTelegram.New(userTelegram.Dependencies{
		Usecase: userUsecase,
		Bot:     bot,
	})

	dailyProgramTelegram := dailyprogramTelegram.New(dailyprogramTelegram.Dependencies{
		Usecase: dailyProgramUsecase,
		Cache:   cache,
		Bot:     bot,
	})

	emitGateway.AddTelegramHandlers(map[event.EventType]emitgateway.TelegramHandler{
		core.NotifySmokingTimeEventType: dailyProgramTelegram.NotifySmokingTime,
	})

	bot.AddHandlers(
		append(
			update.NewFormHandlers(ctx, &update.NewFormHandlerParams[userTelegram.StartProgramForm]{
				Type:     update.HandlerTypeKeyboardCallback,
				Key:      userTelegram.StartFormKey,
				Form:     &userTelegram.StartProgramForm{},
				OnSubmit: userTelegramD.NewStartProgramFormFilled(ctx),
				Bot:      bot,
			}),
			update.NewMessageUpdate(
				userTelegram.StartKey,
				userTelegramD.NewStartBotHandler(ctx),
			),
			update.NewMessageUpdate(
				dailyprogramTelegram.SmokedKey,
				dailyProgramTelegram.NewSmokedHandler(ctx),
			),
			update.NewKeyboardCallbackUpdate(
				dailyprogramTelegram.SmokedButtonKey,
				dailyProgramTelegram.NewSmokedHandler(ctx),
			),
			update.NewKeyboardCallbackUpdate(
				dailyprogramTelegram.SkippedButtonKey,
				dailyProgramTelegram.NewSkippedHandler(ctx),
			),
		)...,
	)

	errChan := make(chan error)
	go func() {
		for err := range errChan {
			fmt.Printf("Received error: %v\n", err)
		}
	}()

	asynqServer := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisUrl},
		asynq.Config{
			Concurrency: 10,
		},
	)

	asynqHandler := asynqhandler.New(asynqhandler.Dependencies{
		DailyProgramUsecase: dailyProgramUsecase,
		EventEmitter:        eventEmitter,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(core.StartProgramQueueEventType.String(),
		asynqHandler.HandleStartProgramQueue,
	)
	mux.HandleFunc(core.ScheduleSmokingTimeQueueEventType.String(),
		asynqHandler.HandleSmokingTimeQueue,
	)
	mux.HandleFunc(core.CalculateNextProgramQueueEventType.String(),
		asynqHandler.HandleCalculateNextProgramQueue,
	)

	go func() {
		handleErr(asynqServer.Run(mux))
	}()

	bot.Updates(ctx, errChan)
}
