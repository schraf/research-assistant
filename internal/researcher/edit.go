package researcher

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

const (
	EditSystemPrompt = `
		You are an expert editor. Your role is to review content and 
		ensure that it takes a neutral stance and is clearly written.
		Remove any Markdown, LaTeX, HTML tags, or any escape characters. The
		content should not include any headings. 
		`
)

func (p *Pipeline) Edit(ctx context.Context, in <-chan Section, out chan<- Section, concurrency int) error {
	defer close(out)

	group, ctx := errgroup.WithContext(ctx)

	for i := 0; i < concurrency; i++ {
		group.Go(func() error {
			for section := range in {
				body, err := p.assistant.Ask(ctx, EditSystemPrompt, section.Body)
				if err != nil {
					return fmt.Errorf("edit error: assistant ask: %w", err)
				}

				section.Body = *body

				slog.Info("edited_section",
					slog.String("section", section.Title),
					slog.Int("length", len(section.Body)),
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
