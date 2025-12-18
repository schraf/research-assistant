package researcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/schraf/assistant/pkg/models"
)

func Compose(ctx context.Context, sections []Section) (*models.Document, error) {
	doc := models.Document{}

	prompt := "Generate a short document title based on its content."

	for _, section := range sections {
		prompt += "\n\n" + section.Body
		doc.AddSection(section.Title, section.Body)
	}

	schema := map[string]any{
		"type":        "string",
		"description": "document title",
	}

	responseJson, err := structuredAsk(ctx, "Create a document title", prompt, schema)
	if err != nil {
		return nil, fmt.Errorf("generate document title error: assistant ask: %w", err)
	}

	if err := json.Unmarshal(responseJson, &doc.Title); err != nil {
		return nil, fmt.Errorf("generate document title error: unmarshal json: %w", err)
	}

	slog.Info("titled_document",
		slog.String("title", doc.Title),
	)

	return &doc, nil
}
