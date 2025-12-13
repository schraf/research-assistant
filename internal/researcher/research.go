package researcher

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/sync/errgroup"
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

func (p *Pipeline) Research(ctx context.Context, in <-chan Section, out chan<- Section, concurrency int) error {
	defer close(out)

	group, ctx := errgroup.WithContext(ctx)

	for i := 0; i < concurrency; i++ {
		group.Go(func() error {
			for section := range in {
				for i := 0; i < ResearchLoops; i++ {
					prompt, err := BuildPrompt(ResearchPrompt, PromptArgs{
						"Topic":       section.Topic + " - " + section.Title,
						"Information": section.Summary + "\n" + section.Research,
					})
					if err != nil {
						return fmt.Errorf("research error: %w", err)
					}

					research, err := p.assistant.Ask(ctx, ResearchSystemPrompt, *prompt)
					if err != nil {
						return fmt.Errorf("research error: assistant ask: %w", err)
					}

					section.Research = *research
				}

				slog.Info("resarched_section",
					slog.String("section", section.Title),
					slog.Int("length", len(section.Research)),
				)

				select {
				case <-ctx.Done():
					return ctx.Err()
				case out <- section:
				}
			}

			return nil
		})
	}

	return group.Wait()
}
