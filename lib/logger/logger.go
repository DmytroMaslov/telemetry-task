package logger

import (
	"log/slog"
	"os"
)

var _ = InnitLogger()

func InnitLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(),
	}))
	slog.SetDefault(logger)
	return logger
}

func LoggerWithPrefix(prefix string) *slog.Logger {
	return slog.Default().With("prefix", prefix)
}

func getLogLevel() slog.Level {
	level, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		return slog.LevelDebug
	}

	switch level {
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	default:
		return slog.LevelDebug
	}
}
