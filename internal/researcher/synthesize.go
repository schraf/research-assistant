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
		results of researched information and synthesize in into a set 
		of paragraphs.
		`

	SynthesizePrompt = `
		# Topic
		{{.Topic}}

		## Researched Information
		{{.Knowledge}}

		# Goal 
		1. Create a detailed series of 4 to 8 paragaphs about the given topic using the information that has been researched and provided. 
		2. Each paragraph should be between 4 to 8 sentences in length.
		`
)

func (p *Pipeline) SynthesizeKnowledge(ctx context.Context, topic string, in <-chan string, out chan<- models.DocumentSection) error {
	defer close(out)

	for knowledge := range in {
		slog.Info("synthesizing_knowledge",
			slog.String("topic", topic),
			slog.String("knowledge", knowledge),
		)

		prompt, err := BuildPrompt(SynthesizePrompt, PromptArgs{
			"Topic":     topic,
			"Knowledge": knowledge,
		})
		if err != nil {
			return fmt.Errorf("synthesize knowledge error: %w", err)
		}

		schema := map[string]any{
			"type":        "array",
			"description": "list of paragraphs",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "name or summary of the paragraph",
					},
					"body": map[string]any{
						"type":        "string",
						"description": "entire body of the paragraph",
					},
				},
				"required": []string{"name", "body"},
			},
		}

		response, err := p.assistant.StructuredAsk(ctx, SynthesizeSystemPrompt, *prompt, schema)
		if err != nil {
			return fmt.Errorf("synthesize knowledge error: assistant structured ask: %w", err)
		}

		var paragraphs []struct {
			Name string `json:"name"`
			Body string `json:"body"`
		}

		if err := json.Unmarshal(response, &paragraphs); err != nil {
			return fmt.Errorf("synthesize knowledge error: json unmarshal: %w", err)
		}

		length := 0
		summaries := []string{}

		for _, paragraph := range paragraphs {
			length += len(paragraph.Body)
			summaries = append(summaries, paragraph.Name)
		}

		slog.Info("synthesized_knowledge",
			slog.String("topic", topic),
			slog.Int("length", length),
			slog.Any("summary", summaries),
		)

		section := models.DocumentSection{
			Title: topic,
		}

		for _, paragraph := range paragraphs {
			section.Paragraphs = append(section.Paragraphs, paragraph.Body)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- section:
		}
	}

	return nil
}
