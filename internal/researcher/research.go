package researcher

import (
	"context"
	"fmt"
	"log/slog"
)

const (
	ResearchSystemPrompt = `
		You are an expert Researcher. You evaluate information
		and use web searches to gather more facts around
		the topic.
		`

	ResearchPrompt = `
		## Topic
		{{.Topic}}

		## Information
		{{.Information}}

		## Goal
		Analyze all of the information and perform web searches
		to gather more facts around the topic that require more
		information. Return a complete summary of all of the
		information gathered.
		`
)

const ResearchLoops = 3

func Research(ctx context.Context, section Section) (*Section, error) {
	for i := 0; i < ResearchLoops; i++ {
		prompt, err := BuildPrompt(ResearchPrompt, PromptArgs{
			"Topic":       section.Topic + " - " + section.Title,
			"Information": section.Summary + "\n" + section.Research,
		})
		if err != nil {
			return nil, fmt.Errorf("research error: %w", err)
		}

		research, err := ask(ctx, ResearchSystemPrompt, *prompt)
		if err != nil {
			return nil, fmt.Errorf("research error: assistant ask: %w", err)
		}

		section.Research = *research
	}

	slog.Info("resarched_section",
		slog.String("section", section.Title),
		slog.Int("length", len(section.Research)),
	)

	return &section, nil
}
