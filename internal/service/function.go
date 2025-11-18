package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
	"github.com/schraf/research-assistant/internal/auth"
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

	// Validate bearer token
	authHeader := r.Header.Get("Authorization")
	if !auth.ValidateToken(authHeader) {
		logger.Warn("unauthorized_request",
			slog.String("auth_header_present", fmt.Sprintf("%v", authHeader != "")),
		)

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract topic from query parameter
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		logger.Error("missing_topic_parameter")

		http.Error(w, "missing required query parameter: topic", http.StatusBadRequest)
		return
	}

	logger.Info("received_research_request",
		slog.String("topic", topic),
	)

	// Execute Cloud Run Job with research request
	if err := executeResearchJob(ctx, logger, requestId, topic); err != nil {
		logger.Error("failed_executing_research_job",
			slog.String("error", err.Error()),
			slog.String("topic", topic),
		)

		http.Error(w, requestErrorMessage(requestId), http.StatusInternalServerError)
		return
	}

	logger.Info("research_job_executed",
		slog.String("topic", topic),
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
