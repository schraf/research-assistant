package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/schraf/research-assistant/internal/service"
	"github.com/schraf/research-assistant/internal/utils"
)

func main() {
	ctx := context.Background()

	if err := utils.LoadEnv(".env"); err != nil {
		slog.Error("load_env_failed",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	if err := utils.SetupLogger("logs/research.log", slog.LevelDebug); err != nil {
		slog.Error("failed_log_setup",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	if err := service.StartServer(ctx); err != nil {
		slog.Error("failed_starting_service",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}
}
