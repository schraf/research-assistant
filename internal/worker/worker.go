package worker

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/schraf/research-assistant/internal/gemini"
	"github.com/schraf/research-assistant/internal/mail"
	"github.com/schraf/research-assistant/internal/models"
	"github.com/schraf/research-assistant/internal/researcher"
	"github.com/schraf/research-assistant/internal/telegraph"
)

func ProcessResearchJob(ctx context.Context, logger *slog.Logger, request models.ResearchRequest) error {
	// Create research report
	report, err := createResearchReport(ctx, logger, request)
	if err != nil {
		return fmt.Errorf("failed researching topic: %w", err)
	}

	// Post research report to Telegraph
	url, err := postResearchReport(ctx, *report)
	if err != nil {
		return fmt.Errorf("failed posting research report: %w", err)
	}

	logger.Info("research_report_posted",
		slog.String("url", *url),
		slog.String("report_title", report.Title),
	)

	// Send email notification
	if err := mail.SendEmail(ctx, logger, report.Title, *url); err != nil {
		return fmt.Errorf("failed sending email: %w", err)
	}

	return nil
}

// createResearchReport creates a research report for the given topic
func createResearchReport(ctx context.Context, logger *slog.Logger, request models.ResearchRequest) (*researcher.ResearchReport, error) {
	client, err := gemini.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	report, err := researcher.ResearchTopic(ctx, logger, client, request.Topic, request.Mode, request.Depth)
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
