package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/schraf/research-assistant/internal/models"
	"github.com/schraf/research-assistant/internal/utils"
	"github.com/schraf/research-assistant/internal/worker"
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

	// Read research request from environment variable
	encodedRequest := os.Getenv("RESEARCH_REQUEST")
	if encodedRequest == "" {
		logger.Error("no_request_data",
			slog.String("error", "no request data found in RESEARCH_REQUEST environment variable"),
		)

		os.Exit(1)
	}

	requestJson, err := base64.StdEncoding.DecodeString(encodedRequest)
	if err != nil {
		logger.Error("failed_decoding_request",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	var request models.ResearchRequest
	if err := json.Unmarshal(requestJson, &request); err != nil {
		logger.Error("failed_parsing_request_json",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	logger = logger.With(slog.String("request_id", request.RequestId))

	// Validate request fields
	if request.RequestId == "" {
		logger.Error("invalid_request_id",
			slog.String("error", "missing request id"),
		)

		os.Exit(1)
	}

	if request.Topic == "" {
		logger.Error("invalid_request_id",
			slog.String("error", "missing topic"),
		)

		os.Exit(1)
	}

	if err := request.Depth.Validate(); err != nil {
		logger.Error("invalid_request_id",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	if err := request.Mode.Validate(); err != nil {
		logger.Error("invalid_request_id",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	// Process the research job
	if err := worker.ProcessResearchJob(ctx, logger, request); err != nil {
		logger.Error("job_failed",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	logger.Info("job_completed_successfully")
	os.Exit(0)
}
