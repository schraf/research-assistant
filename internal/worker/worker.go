package worker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/schraf/research-assistant/internal/gemini"
	"github.com/schraf/research-assistant/internal/mail"
	"github.com/schraf/research-assistant/internal/models"
	"github.com/schraf/research-assistant/internal/researcher"
	"github.com/schraf/research-assistant/internal/telegraph"
)

// ProcessResearchJob processes a Pub/Sub message containing a research request
// It expects the message data to be base64-encoded JSON matching ResearchRequest
func ProcessResearchJob(ctx context.Context, logger *slog.Logger, messageData string) error {
	// Decode base64 message data
	decoded, err := base64.StdEncoding.DecodeString(messageData)
	if err != nil {
		return fmt.Errorf("failed to decode base64 message data: %w", err)
	}

	// Parse ResearchRequest from JSON
	var req models.ResearchRequest
	if err := json.Unmarshal(decoded, &req); err != nil {
		return fmt.Errorf("failed to parse research request: %w", err)
	}

	// Validate required fields
	if req.RequestId == "" {
		return fmt.Errorf("request_id is required")
	}
	if req.Topic == "" {
		return fmt.Errorf("topic is required")
	}

	// Create logger with request_id for traceability
	jobLogger := logger.With(slog.String("request_id", req.RequestId))

	jobLogger.Info("processing_research_job",
		slog.String("topic", req.Topic),
	)

	// Create research report
	report, err := createResearchReport(ctx, jobLogger, req.Topic)
	if err != nil {
		jobLogger.Error("failed_researching_topic",
			slog.String("error", err.Error()),
			slog.String("topic", req.Topic),
		)
		return fmt.Errorf("failed researching topic: %w", err)
	}

	if report == nil {
		jobLogger.Error("research_report_is_nil",
			slog.String("topic", req.Topic),
		)
		return fmt.Errorf("research report is nil")
	}

	// Post research report to Telegraph
	url, err := postResearchReport(ctx, *report)
	if err != nil {
		jobLogger.Error("failed_posting_research_report",
			slog.String("error", err.Error()),
			slog.String("topic", req.Topic),
			slog.String("report_title", report.Title),
		)
		return fmt.Errorf("failed posting research report: %w", err)
	}

	if url == nil || *url == "" {
		jobLogger.Error("telegraph_url_is_empty",
			slog.String("topic", req.Topic),
			slog.String("report_title", report.Title),
		)
		return fmt.Errorf("telegraph URL is empty")
	}

	jobLogger.Info("research_report_posted",
		slog.String("url", *url),
		slog.String("topic", req.Topic),
		slog.String("report_title", report.Title),
	)

	// Send email notification
	if err := mail.SendEmail(ctx, jobLogger, report.Title, *url); err != nil {
		jobLogger.Error("failed_sending_email",
			slog.String("error", err.Error()),
			slog.String("topic", req.Topic),
			slog.String("report_title", report.Title),
			slog.String("url", *url),
		)
		return fmt.Errorf("failed sending email: %w", err)
	}

	jobLogger.Info("research_job_completed",
		slog.String("topic", req.Topic),
		slog.String("report_title", report.Title),
		slog.String("url", *url),
	)

	return nil
}

// createResearchReport creates a research report for the given topic
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

// postResearchReport posts a research report to Telegraph and returns the URL
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

	authorName := os.Getenv("TELEGRAPH_AUTHOR_NAME")
	returnContent := false

	pageRequest := telegraph.CreatePageRequest{
		AccessToken:   apiToken,
		Title:         report.Title,
		AuthorName:    &authorName,
		Content:       content,
		ReturnContent: &returnContent,
	}

	page, err := client.CreatePage(ctx, pageRequest)
	if err != nil {
		return nil, err
	}

	return &page.URL, nil
}
