package tlog

import (
	"log/slog"
	"os"
	"shantaram/pkg/build"
	"shantaram/pkg/config"

	slogmulti "github.com/samber/slog-multi"
)

func Init(cfg *config.Config) error {
	logHandlers := []slog.Handler{slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	})}

	multiHandler := slogmulti.Fanout(logHandlers...)
	ctxHandler := &contextHandler{multiHandler}

	logger := slog.New(ctxHandler).With(
		slog.String("app", "api"),
		slog.String("app_tag", build.Tag),
	)
	slog.SetDefault(logger)

	return nil
}
