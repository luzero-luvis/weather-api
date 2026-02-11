package logger

import (
	"log/slog"
	"os"
)

func setUp(env string) *slog.Logger {
	var handler slog.Handler

	if env == "development" || env == "" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}
