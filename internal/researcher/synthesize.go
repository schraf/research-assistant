package researcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/schraf/assistant/pkg/models"
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

func SynthesizeReport(ctx context.Context, assistant models.Assistant, topic string, results []ResearchResult, depth ResearchDepth) (*models.Document, error) {
	slog.InfoContext(ctx, "synthesizing_report")

	prompt, err := BuildPrompt(SynthesizePrompt, PromptArgs{
		"ResearchTopic":   topic,
		"ResearchResults": results,
	})
	if err != nil {
		return nil, fmt.Errorf("failed building research report prompt: %w", err)
	}

	response, err := assistant.StructuredAsk(ctx, SynthesizeSystemPrompt, *prompt, SynthesizeReportSchema())
	if err != nil {
		return nil, fmt.Errorf("failed syntesizing research report: %w", err)
	}

	var doc models.Document

	if err := json.Unmarshal(response, &doc); err != nil {
		return nil, fmt.Errorf("failed parsing research report: %w", err)
	}

	return &doc, nil
}

func SynthesizeReportSchema() map[string]any {
	return map[string]any{
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
						"title": map[string]any{
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
						"title",
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
