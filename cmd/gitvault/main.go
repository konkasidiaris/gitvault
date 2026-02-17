package main

import (
	"log/slog"
	"os"

	"github.com/konkasidiaris/gitvault/internal/logging"
	"github.com/konkasidiaris/gitvault/internal/sync"
)

func main() {
	logging.InitializeLogger()

	if err := sync.Run(); err != nil {
		slog.Error("sync failed", "error", err)
		os.Exit(1)
	}
}
