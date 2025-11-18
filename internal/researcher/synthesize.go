package researcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/schraf/research-assistant/internal/models"
)

const (
	SynthesizeSystemPrompt = `
		You are an expert Research Writer. Your sole task is to take the
		results of researched information and synthesize a formal research
		report.  
		`

	SynthesizePrompt = `
		# Research Topic
		{{.ResearchTopic}}

		# Researched Information

		{{range $_, $item := .ResearchResults}}
		### {{$item.Topic}}

		{{range $_, $knowledge := $item.Knowledge}}
        #### {{$knowledge.Topic}}
		{{$knowledge.Information}}

		{{end}}
		{{end}}

		# Goal 
		Create a structured and detailed and formal research report of the
		given topic using the information that has been researched and
		provided. The report should have a title. The report should be composed
		of a series of sections where each section is a list of paragraphs. The
		paragraphs should be in in the english language and follow correct
		grammar.  
		`
)

type ReportSection struct {
	SectionTitle string   `json:"section_title"`
	Paragraphs   []string `json:"paragraphs"`
}

type ResearchReport struct {
	Title    string          `json:"title"`
	Sections []ReportSection `json:"sections"`
}

func SynthesizeReport(ctx context.Context, logger *slog.Logger, resources models.Resources, topic string, results []ResearchResult) (*ResearchReport, error) {
	logger.InfoContext(ctx, "synthesizing_report")

	prompt, err := BuildPrompt(SynthesizePrompt, PromptArgs{
		"ResearchTopic":   topic,
		"ResearchResults": results,
	})
	if err != nil {
		return nil, fmt.Errorf("failed building research report prompt: %w", err)
	}

	response, err := resources.StructuredAsk(ctx, SynthesizeSystemPrompt, *prompt, SynthesizeReportSchema())
	if err != nil {
		return nil, fmt.Errorf("failed syntesizing research report: %w", err)
	}

	var report ResearchReport

	if err := json.Unmarshal(response, &report); err != nil {
		return nil, fmt.Errorf("failed parsing research report: %w", err)
	}

	return &report, nil
}

func SynthesizeReportSchema() models.Schema {
	return models.Schema{
		"type": "object",
		"properties": map[string]any{
			"title": map[string]any{
				"type":        "string",
				"description": "The research title for the report",
			},
			"sections": map[string]any{
				"type":        "array",
				"description": "A list of sections for the report",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"section_title": map[string]any{
							"type":        "string",
							"description": "A title for the section in the report",
						},
						"paragraphs": map[string]any{
							"type":        "array",
							"description": "A list of separate paragraphs for the section of the report",
							"items": map[string]any{
								"type": "string",
							},
						},
					},
					"required": []string{
						"section_title",
						"paragraphs",
					},
				},
			},
		},
		"required": []string{
			"title",
			"sections",
		},
	}
}
