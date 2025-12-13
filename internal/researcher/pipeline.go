package researcher

import (
	"context"

	"github.com/schraf/assistant/pkg/models"
	"golang.org/x/sync/errgroup"
)

type Pipeline struct {
	assistant models.Assistant
}

func NewPipeline(assistant models.Assistant) *Pipeline {
	return &Pipeline{
		assistant: assistant,
	}
}

func (p *Pipeline) Exec(ctx context.Context, topic string) (*models.Document, error) {
	plan := make(chan Section)
	research := make(chan Section, 6)
	synthesis := make(chan Section, 6)
	edited := make(chan Section, 6)
	aggregated := make(chan models.Document, 1)
	titled := make(chan models.Document, 1)

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return p.Plan(ctx, topic, plan)
	})

	group.Go(func() error {
		return p.Research(ctx, plan, research, 3)
	})

	group.Go(func() error {
		return p.Synthesize(ctx, research, synthesis, 3)
	})

	group.Go(func() error {
		return p.Edit(ctx, synthesis, edited, 3)
	})

	group.Go(func() error {
		return p.Aggregate(ctx, edited, aggregated)
	})

	group.Go(func() error {
		return p.Title(ctx, aggregated, titled)
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	doc := <-titled

	return &doc, nil
}
