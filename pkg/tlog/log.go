package tlog

import (
	"log/slog"
	"os"
	"shantaram/pkg/build"
	"shantaram/pkg/config"
	"shantaram/pkg/telemetry"

	slogmulti "github.com/samber/slog-multi"
	slogtelegram "github.com/samber/slog-telegram/v2"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func Init(cfg *config.Config, tel *telemetry.Telemetry) error {
	logHandlers := []slog.Handler{slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	})}

	if cfg.Telemetry.Enabled {
		logHandlers = append(logHandlers, otelslog.NewHandler(cfg.ServiceName,
			otelslog.WithSource(true),
			otelslog.WithLoggerProvider(tel.LogProvider),
		))
	}

	if cfg.Log.Telegram.Token != "" {
		logHandlers = append(logHandlers, slogtelegram.Option{
			Level:     slog.LevelError,
			Token:     cfg.Log.Telegram.Token,
			Username:  cfg.Log.Telegram.ChatID,
			AddSource: true,
		}.NewTelegramHandler())
	}

	multiHandler := slogmulti.Fanout(logHandlers...)
	ctxHandler := &contextHandler{multiHandler}

	logger := slog.New(ctxHandler).With(
		slog.String(string(semconv.ServiceNameKey), cfg.ServiceName),
		slog.String(string(semconv.ServiceVersionKey), build.Tag),
	)
	slog.SetDefault(logger)

	return nil
}
