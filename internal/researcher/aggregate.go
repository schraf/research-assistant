package researcher

import (
	"context"

	"github.com/schraf/assistant/pkg/models"
)

func (p *Pipeline) Aggregate(ctx context.Context, in <-chan Section, out chan<- models.Document) error {
	defer close(out)

	doc := models.Document{}

	for section := range in {
		doc.AddSection(section.Title, section.Body)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case out <- doc:
	}

	return nil
}
