package researcher

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

const (
	ResearchSystemPrompt = `
		You are an expert Researcher. Your sole task is
		to search the web to gather information to answer
		a question.
		`

	ResearchPrompt = `
		## Question
		{{.Question}}
		
		## Goal 
		Search the web and gather information that will answer
		the provided question. You do not need to respond with 
		the answer, but rather provide plenty of information 
		around the topic that can be used to answer the question.
		Please provide the information in a well formatted structure.
		`
)

func (p *Pipeline) ResearchQuestion(ctx context.Context, in <-chan string, out chan<- string, concurrency int) error {
	defer close(out)

	group, ctx := errgroup.WithContext(ctx)

	for i := 0; i < concurrency; i++ {
		group.Go(func() error {
			for question := range in {
				slog.Info("researching_question",
					slog.String("question", question),
				)

				prompt, err := BuildPrompt(ResearchPrompt, PromptArgs{
					"Question": question,
				})
				if err != nil {
					return fmt.Errorf("research question error: %w", err)
				}

				information, err := p.assistant.Ask(ctx, ResearchSystemPrompt, *prompt)
				if err != nil {
					return fmt.Errorf("research question error: assistant ask: %w", err)
				}

				slog.Info("resarched_question",
					slog.String("question", question),
					slog.String("information", *information),
				)

				select {
				case <-ctx.Done():
					return ctx.Err()
				case out <- *information:
				}
			}

			return nil
		})
	}

	return group.Wait()
}
