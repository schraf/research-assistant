package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
	"github.com/schraf/research-assistant/internal/models"
)

func init() {
	functions.HTTP("research", research)
}

func requestErrorMessage(requestId string) string {
	return "An internal error has occurred. (" + requestId + ")"
}

func research(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := uuid.NewString()
	logger := slog.Default().With(slog.String("request_id", requestId))

	// Decode JSON payload
	var payload struct {
		Topic string               `json:"topic"`
		Mode  models.ResourceMode  `json:"resource_mode"`
		Depth models.ResearchDepth `json:"research_depth"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Error("invalid_json_payload",
			slog.String("error", err.Error()),
		)

		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	if payload.Topic == "" {
		logger.Error("missing_topic_field")

		http.Error(w, "missing required field: topic", http.StatusBadRequest)
		return
	}

	// Build research request model
	req := models.ResearchRequest{
		RequestId: requestId,
		Mode:      payload.Mode,
		Depth:     payload.Depth,
		Topic:     payload.Topic,
	}

	logger.Info("received_research_request",
		slog.String("topic", req.Topic),
		slog.Int("resource_mode", int(req.Mode)),
		slog.Int("research_depth", int(req.Depth)),
	)

	// Execute Cloud Run Job with research request
	if err := executeResearchJob(ctx, logger, req); err != nil {
		logger.Error("failed_executing_research_job",
			slog.String("error", err.Error()),
			slog.String("topic", req.Topic),
		)

		http.Error(w, requestErrorMessage(requestId), http.StatusInternalServerError)
		return
	}

	logger.Info("research_job_executed",
		slog.String("topic", req.Topic),
	)

	// Return success response
	response := map[string]interface{}{
		"success":   true,
		"request_id": requestId,
		"message":   "Research request queued",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("failed_encoding_response",
			slog.String("error", err.Error()),
		)
	}
}
