package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/schraf/gemini-email/internal/gemini"
	"github.com/schraf/gemini-email/internal/researcher"
	"github.com/schraf/gemini-email/internal/telegraph"
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

func CreateResearchReport(ctx context.Context, topic string) (*researcher.ResearchReport, error) {
	client, err := gemini.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	report, err := researcher.ResearchTopic(ctx, client, Topic)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func PostResearchReport(ctx context.Context, report researcher.ResearchReport) (*string, error) {
	client := telegraph.NewDefaultClient()
	content := telegraph.Nodes{}
	apiToken := os.Getenv("TELEGRAPH_API_KEY")

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

	report, err := CreateResearchReport(ctx, Topic)
	if err != nil {
		slog.Error("failed_researching_topic",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	url, err := PostResearchReport(ctx, *report)
	if err != nil {
		slog.Error("failed_posting_research_repo",
			slog.String("error", err.Error()),
		)
	}

	slog.Info("research_report_posted",
		slog.String("url", *url),
	)
}
