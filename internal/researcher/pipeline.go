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
	subtopics := make(chan string)
	sections := make(chan models.DocumentSection)
	out := make(chan models.Document, 1)

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return p.CreateSubtopics(ctx, topic, subtopics)
	})

	group.Go(func() error {
		return p.CreateDocumentSection(ctx, subtopics, sections, 10)
	})

	group.Go(func() error {
		defer close(out)
		doc, err := p.EditDocument(ctx, sections)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- *doc:
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	doc := <-out

	return &doc, nil
}
