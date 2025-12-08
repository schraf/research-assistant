package researcher

import (
	"context"
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
		Create a detailed series of paragaphs about the given topic using the 
		information that has been researched and provided. The paragraphs should 
		be in in the english language and follow correct grammar. Do not include
		any markdown, LaTeX, HTML tags or any special characters.
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

		response, err := p.assistant.Ask(ctx, SynthesizeSystemPrompt, *prompt)
		if err != nil {
			return fmt.Errorf("synthesize knowledge error: assistant ask: %w", err)
		}

		paragraphs, err := GenerateList(ctx, p.assistant, *response)
		if err != nil {
			return fmt.Errorf("synthesize knowledge error: %w", err)
		}

		length := 0

		for _, paragraph := range paragraphs {
			length += len(paragraph)
		}

		slog.Info("synthesized_knowledge",
			slog.String("topic", topic),
			slog.Int("length", length),
			slog.Any("paragraphs", paragraphs),
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- models.DocumentSection{Title: topic, Paragraphs: paragraphs}:
		}
	}

	return nil
}
