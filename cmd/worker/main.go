package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/schraf/research-assistant/internal/worker"
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

	if err := utils.SetupLogger("logs/worker.log", slog.LevelDebug); err != nil {
		slog.Error("failed_log_setup",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	logger := slog.Default()

	// Read CloudEvent from stdin (Cloud Run Jobs pass events via stdin)
	var eventData []byte
	var err error

	// Try to read from stdin first (for Cloud Run Jobs)
	if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		// stdin is a pipe or file
		eventData, err = io.ReadAll(os.Stdin)
		if err != nil {
			logger.Error("failed_reading_stdin",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	} else {
		// stdin is a terminal, try environment variable (for local testing)
		eventDataStr := os.Getenv("CLOUDEVENT_DATA")
		if eventDataStr == "" {
			logger.Error("no_event_data",
				slog.String("error", "no CloudEvent data found in stdin or CLOUDEVENT_DATA environment variable"),
			)
			os.Exit(1)
		}
		eventData = []byte(eventDataStr)
	}

	if len(eventData) == 0 {
		logger.Error("empty_event_data")
		os.Exit(1)
	}

	// Parse CloudEvent JSON
	var cloudEvent map[string]interface{}
	if err := json.Unmarshal(eventData, &cloudEvent); err != nil {
		logger.Error("failed_parsing_cloudevent",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Extract Pub/Sub message data
	data, ok := cloudEvent["data"].(map[string]interface{})
	if !ok {
		logger.Error("invalid_cloudevent_data",
			slog.String("error", "cloudEvent.data is not a map"),
		)
		os.Exit(1)
	}

	message, ok := data["message"].(map[string]interface{})
	if !ok {
		logger.Error("invalid_pubsub_message",
			slog.String("error", "data.message is not a map"),
		)
		os.Exit(1)
	}

	messageData, ok := message["data"].(string)
	if !ok {
		logger.Error("invalid_message_data",
			slog.String("error", "message.data is not a string"),
		)
		os.Exit(1)
	}

	// Process the research job
	if err := worker.ProcessResearchJob(ctx, logger, messageData); err != nil {
		logger.Error("job_failed",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	logger.Info("job_completed_successfully")
	os.Exit(0)
}
