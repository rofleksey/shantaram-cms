package telegram

import (
	"context"
	"fmt"
	"shantaram/pkg/config"
	"shantaram/pkg/database"
	"shantaram/pkg/telemetry"

	tgBot "github.com/go-telegram/bot"
	"github.com/samber/do"
)

var serviceName = "telegram"

type Service struct {
	appCtx  context.Context
	cfg     *config.Config
	queries *database.Queries
	tracing *telemetry.Tracing
	bot     *tgBot.Bot
}

func New(di *do.Injector) (*Service, error) {
	cfg := do.MustInvoke[*config.Config](di)

	bot, err := tgBot.New(cfg.Telegram.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &Service{
		appCtx:  do.MustInvoke[context.Context](di),
		cfg:     cfg,
		queries: do.MustInvoke[*database.Queries](di),
		tracing: do.MustInvoke[*telemetry.Tracing](di),
		bot:     bot,
	}, nil
}

func (s *Service) Notify(msg string) {
	for _, id := range s.cfg.Telegram.ChatIds {
		_, _ = s.bot.SendMessage(s.appCtx, &tgBot.SendMessageParams{
			ChatID: id,
			Text:   msg,
		})
	}
}
