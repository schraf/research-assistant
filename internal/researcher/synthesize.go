package researcher

import (
	"context"
	"fmt"
	"log/slog"
)

const (
	SynthesizeSystemPrompt = `
		You are an expert Research Writer. Your sole task is to take the
		results of researched information and synthesize in into a set 
		of 3 to 5 paragraphs and each paragraph should have at least 4
		sentences. Do not include any headings or Markdown, HTML, LaTeX,
		or escape characters. The paragraphs should be written in clean,
		neutral, report-style English.
		`
)

func Synthesize(ctx context.Context, section Section) (*Section, error) {
	body, err := ask(ctx, SynthesizeSystemPrompt, section.Research)
	if err != nil {
		return nil, fmt.Errorf("synthesize error: assistant ask: %w", err)
	}

	section.Body = *body

	slog.Info("synthesized_section",
		slog.String("section", section.Title),
		slog.Int("body", len(section.Body)),
	)

	return &section, nil
}
