package researcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/schraf/assistant/pkg/models"
)

func (p *Pipeline) Title(ctx context.Context, in <-chan models.Document, out chan<- models.Document) error {
	defer close(out)

	doc := <-in

	prompt := "Generate a short document title based on its content.\n"

	for _, section := range doc.Sections {
		prompt += "\n\n" + strings.Join(section.Paragraphs, "\n")
	}

	schema := map[string]any{
		"type":        "string",
		"description": "document title",
	}

	responseJson, err := p.assistant.StructuredAsk(ctx, "Create a document title", prompt, schema)
	if err != nil {
		return fmt.Errorf("generate document title error: assistant ask: %w", err)
	}

	if err := json.Unmarshal(responseJson, &doc.Title); err != nil {
		return fmt.Errorf("generate document title error: unmarshal json: %w", err)
	}

	slog.Info("titled_document",
		slog.String("title", doc.Title),
	)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case out <- doc:
	}

	return nil
}
