package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sousair/gocore/pkg/cache"
	"github.com/sousair/stop-smoking/internal/core"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
)

func (h TelegramHandler) NotifySmokingTime(ctx context.Context, data []byte) error {
	var req *core.SmokingSchedule
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	cutOff := time.Now().Local().Add(time.Second * 10)

	if req.NextSmoke.After(cutOff) {
		return nil
	}

	msg := message.NewSimpleMessage(&message.Params{
		Recipient: req.Identifier,
		Content: `🚬It's time to smoke! 🚬
You have 5 minutes to smoke or skip your cigarette. After that, will be consider a failure
If you are skipping, we will consider if you want to smoke later.
If you are going to smoke or already smoked, please click the button below.`,
		Bot: h.deps.Bot,
	})

	if err := h.deps.Cache.Set(ctx,
		fmt.Sprintf(CanSmokeKeyPattern, req.Identifier),
		"true",
		cache.WithTTL(time.Minute*5),
	); err != nil {
		return err
	}

	if _, err := msg.Send(ctx, message.WithMessageButtons(message.KeyboardRow{
		{"Smoked", SmokedButtonKey},
		{"Skipped", SkippedButtonKey},
	})); err != nil {
		return err
	}

	return nil
}
