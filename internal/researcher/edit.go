package researcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/schraf/assistant/pkg/models"
)

const (
	EditSystemPrompt = `
		You are an expert research report editor. Your role is to review research
		reports and ensure they flow cohesively, with no repetition of content
		between sections. Each section should build upon previous sections without
		reiterating the same information. The report should read as a unified,
		well-structured document where each section contributes unique value and
		transitions smoothly to the next. As a research report editor, you have
		extensive experience in academic and professional research writing, and you
		know how to refine reports to eliminate redundancy while maintaining clarity
		and coherence.
		`

	EditPrompt = `
		# Research Report to Edit

		## Title
		{{.Title}}

		## Sections

		{{range $index, $section := .Sections}}
		### Section {{$index}}: {{$section.Title}}

		{{range $_, $paragraph := $section.Paragraphs}}
		{{$paragraph}}

		{{end}}
		{{end}}

		# Goal
		Review and edit this research report to ensure:
		1. The report flows cohesively from section to section
		2. There is no repetition of content between sections
		3. Each section builds upon previous sections without reiterating information
		4. Transitions between sections are smooth and logical
		5. The overall report reads as a unified, well-structured document

		Maintain the same structure (title and sections with paragraphs) but refine
		the content to eliminate redundancy and improve flow. If information is
		repeated across sections, consolidate it appropriately or remove redundant
		instances. Ensure each section contributes unique value to the overall report.
		`
)

func EditReport(ctx context.Context, assistant models.Assistant, doc *models.Document) (*models.Document, error) {
	slog.InfoContext(ctx, "editing_report")

	prompt, err := BuildPrompt(EditPrompt, doc)
	if err != nil {
		return nil, fmt.Errorf("failed building edit prompt: %w", err)
	}

	response, err := assistant.StructuredAsk(ctx, EditSystemPrompt, *prompt, EditReportSchema())
	if err != nil {
		return nil, fmt.Errorf("failed editing research report: %w", err)
	}

	var editedDoc models.Document

	if err := json.Unmarshal(response, &editedDoc); err != nil {
		return nil, fmt.Errorf("failed parsing edited research report: %w", err)
	}

	return &editedDoc, nil
}

func EditReportSchema() map[string]any {
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
