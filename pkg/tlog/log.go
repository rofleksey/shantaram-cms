package tlog

import (
	"log/slog"
	"os"
	"shantaram/pkg/build"
	"shantaram/pkg/config"
	"shantaram/pkg/telemetry"

	slogmulti "github.com/samber/slog-multi"
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
		logHandlers = append(logHandlers, otelslog.NewHandler(build.ServiceName,
			otelslog.WithSource(true),
			otelslog.WithLoggerProvider(tel.LogProvider),
		))
	}

	multiHandler := slogmulti.Fanout(logHandlers...)
	ctxHandler := &contextHandler{multiHandler}

	logger := slog.New(ctxHandler).With(
		slog.String(string(semconv.ServiceNameKey), build.ServiceName),
		slog.String(string(semconv.ServiceVersionKey), build.Tag),
	)
	slog.SetDefault(logger)

	return nil
}
