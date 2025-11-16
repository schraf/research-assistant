package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/schraf/gemini-email/internal/gemini"
	"github.com/schraf/gemini-email/internal/researcher"
	"github.com/schraf/gemini-email/internal/utils"
)

const (
	Topic = `
		I would a report of the life and works of H. P. Lovecraft. Please
		include details about this family and personal life. Also include
		information about this most famous works of fiction. Finally include
		any notable details about his legacy.  
		`
)

func main() {
	ctx := context.Background()

	if err := utils.SetupLogger("logs/app.log", slog.LevelDebug); err != nil {
		slog.Error("failed_log_setup",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	client, err := gemini.NewClient(ctx)
	if err != nil {
		slog.Error("failed_creating_gemini_client",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	closer, err := client.EnableLogging("chat.log")
	if err != nil {
		slog.Error("failed_gemini_logging",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}
	defer closer()

	if _, err := researcher.ResearchTopic(ctx, client, Topic); err != nil {
		slog.Error("failed_researching_topic",
			slog.String("topic", Topic),
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}
}
