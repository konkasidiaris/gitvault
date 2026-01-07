package logging

import (
	"log/slog"
	"os"
)

func InitializeLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
}
