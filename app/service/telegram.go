package service

import (
	"context"
	"fmt"
	tgBot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"shantaram-cms/pkg/config"
	"strconv"
)

type Telegram struct {
	ctx context.Context
	cfg *config.Config
	bot *tgBot.Bot
}

func NewTelegram(ctx context.Context, cfg *config.Config) (*Telegram, error) {
	opts := []tgBot.Option{
		tgBot.WithDefaultHandler(func(ctx context.Context, b *tgBot.Bot, update *models.Update) {
			_, _ = b.SendMessage(ctx, &tgBot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Chat ID = " + strconv.FormatInt(update.Message.Chat.ID, 10),
			})
		}),
	}

	bot, err := tgBot.New(cfg.TelegramToken, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	go bot.Start(ctx)

	return &Telegram{
		ctx: ctx,
		cfg: cfg,
		bot: bot,
	}, nil
}

func (s *Telegram) Notify(msg string) {
	for _, id := range s.cfg.TelegramChatIDs {
		_, _ = s.bot.SendMessage(s.ctx, &tgBot.SendMessageParams{
			ChatID: id,
			Text:   msg,
		})
	}
}
