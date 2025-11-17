package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
	"github.com/schraf/research-assistant/internal/gemini"
	"github.com/schraf/research-assistant/internal/mail"
	"github.com/schraf/research-assistant/internal/researcher"
	"github.com/schraf/research-assistant/internal/telegraph"
)

func init() {
	functions.HTTP("research", research)
}

func createResearchReport(ctx context.Context, logger *slog.Logger, topic string) (*researcher.ResearchReport, error) {
	client, err := gemini.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	report, err := researcher.ResearchTopic(ctx, logger, client, topic)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func postResearchReport(ctx context.Context, report researcher.ResearchReport) (*string, error) {
	client := telegraph.NewDefaultClient()
	content := telegraph.Nodes{}
	apiToken := os.Getenv("TELEGRAPH_API_KEY")
	if apiToken == "" {
		return nil, fmt.Errorf("TELEGRAPH_API_KEY environment variable is not set")
	}

	for _, section := range report.Sections {
		//--===  ADD SECTION TITLE
		content = append(content, telegraph.NodeElement{
			Tag: "h3",
			Children: telegraph.Nodes{
				section.SectionTitle,
			},
		})

		//--=== ADD SECTION PARAGRAPHS
		for _, paragraph := range section.Paragraphs {
			content = append(content, telegraph.NodeElement{
				Tag: "p",
				Children: telegraph.Nodes{
					paragraph,
				},
			})
		}
	}

	pageRequest := telegraph.CreatePageRequest{
		AccessToken: apiToken,
		Title:       report.Title,
		Content:     content,
	}

	page, err := client.CreatePage(ctx, pageRequest)
	if err != nil {
		return nil, err
	}

	return &page.URL, nil
}

func requestErrorMessage(requestId string) string {
	return "An internal error has occurred. (" + requestId + ")"
}

func research(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := uuid.NewString()
	logger := slog.With(slog.String("request_id", requestId))

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

	// Create research report
	report, err := createResearchReport(ctx, logger, topic)
	if err != nil {
		logger.Error("failed_researching_topic",
			slog.String("error", err.Error()),
			slog.String("topic", topic),
		)

		http.Error(w, requestErrorMessage(requestId), http.StatusInternalServerError)
		return
	}

	if report == nil {
		logger.Error("research_report_is_nil",
			slog.String("topic", topic),
		)

		http.Error(w, requestErrorMessage(requestId), http.StatusInternalServerError)
		return
	}

	// Post research report to Telegraph
	url, err := postResearchReport(ctx, *report)
	if err != nil {
		logger.Error("failed_posting_research_report",
			slog.String("error", err.Error()),
			slog.String("topic", topic),
			slog.String("report_title", report.Title),
		)

		http.Error(w, requestErrorMessage(requestId), http.StatusInternalServerError)
		return
	}

	if url == nil || *url == "" {
		logger.Error("telegraph_url_is_empty",
			slog.String("topic", topic),
			slog.String("report_title", report.Title),
		)

		http.Error(w, requestErrorMessage(requestId), http.StatusInternalServerError)
		return
	}

	logger.Info("research_report_posted",
		slog.String("url", *url),
		slog.String("topic", topic),
		slog.String("report_title", report.Title),
	)

	// Send email notification (non-critical, log but don't fail the request)
	if err := mail.SendEmail(ctx, logger, report.Title, *url); err != nil {
		logger.Error("failed_sending_email",
			slog.String("error", err.Error()),
			slog.String("topic", topic),
			slog.String("report_title", report.Title),
			slog.String("url", *url),
		)

		http.Error(w, requestErrorMessage(requestId), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"success":   true,
		"report_id": requestId,
		"title":     report.Title,
		"url":       *url,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("failed_encoding_response",
			slog.String("error", err.Error()),
		)
	}
}
