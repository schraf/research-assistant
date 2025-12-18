package researcher

import (
	"context"

	"github.com/schraf/assistant/pkg/models"
	"github.com/schraf/pipeline"
)

func CreateDocument(ctx context.Context, assistant models.Assistant, topic string) (*models.Document, error) {
	ctx = withAssistant(ctx, assistant)

	plan, err := Plan(ctx, topic)
	if err != nil {
		return nil, err
	}

	pipe, ctx := pipeline.WithPipeline(ctx)

	researched := make(chan Section, 6)
	pipeline.ParallelTransform(pipe, 6, Research, plan, researched)

	synthesis := make(chan Section, 6)
	pipeline.ParallelTransform(pipe, 6, Synthesize, researched, synthesis)

	edited := make(chan Section, 6)
	pipeline.ParallelTransform(pipe, 6, Edit, synthesis, edited)

	aggregated := make(chan []Section, 1)
	pipeline.Aggregate(pipe, edited, aggregated)

	composed := make(chan models.Document, 1)
	pipeline.Transform(pipe, Compose, aggregated, composed)

	if err := pipe.Wait(); err != nil {
		return nil, err
	}

	doc := <-composed

	return &doc, nil
}
