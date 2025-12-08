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
		You are an expert editor. Your role is to review documents and ensure that
		each section contributes unique value and transitions smoothly to the next. 
		You do not perform any further research or add new content.
		`

	EditPrompt = `
		# Research Report to Edit

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
		6. Decided on a title for the report
		7. Remove any markdown, LaTeX, HTML tags, or any escape characters.
		`
)

func (p *Pipeline) EditDocument(ctx context.Context, in <-chan models.DocumentSection) (*models.Document, error) {
	aggregated := []models.DocumentSection{}

	for section := range in {
		aggregated = append(aggregated, section)
	}

	slog.Info("editing_document",
		slog.Int("sections_count", len(aggregated)),
	)

	draft := models.Document{
		Sections: aggregated,
	}

	prompt, err := BuildPrompt(EditPrompt, draft)
	if err != nil {
		return nil, fmt.Errorf("edit document error: %w", err)
	}

	response, err := p.assistant.StructuredAsk(ctx, EditSystemPrompt, *prompt, DocumentSchema())
	if err != nil {
		return nil, fmt.Errorf("edit document error: assistant structured ask: %w", err)
	}

	var doc models.Document

	if err := json.Unmarshal(response, &doc); err != nil {
		return nil, fmt.Errorf("edit document error: json unmarshal: %w", err)
	}

	slog.Info("edited_document",
		slog.Int("draft_length", DocumentLength(&draft)),
		slog.Int("final_length", DocumentLength(&doc)),
	)

	return &doc, nil
}

func DocumentSchema() map[string]any {
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
