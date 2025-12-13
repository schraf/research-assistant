package researcher

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/sync/errgroup"
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

func (p *Pipeline) Synthesize(ctx context.Context, in <-chan Section, out chan<- Section, concurrency int) error {
	defer close(out)

	group, ctx := errgroup.WithContext(ctx)

	for i := 0; i < concurrency; i++ {
		group.Go(func() error {
			for section := range in {
				body, err := p.assistant.Ask(ctx, SynthesizeSystemPrompt, section.Research)
				if err != nil {
					return fmt.Errorf("synthesize error: assistant ask: %w", err)
				}

				section.Body = *body

				slog.Info("synthesized_section",
					slog.String("section", section.Title),
					slog.Int("body", len(section.Body)),
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
