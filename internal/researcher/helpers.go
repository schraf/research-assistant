package researcher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/schraf/assistant/pkg/models"
)

const (
	ListPrompt = `
		Given the following list. Break it down into the core elements. 
		Return a structured list that follows the provided schema ensuring
		that each item is in its own element in the list.

		{{.Input}}
		`
)

func GenerateList(ctx context.Context, assistant models.Assistant, input string) ([]string, error) {
	prompt, err := BuildPrompt(ListPrompt, PromptArgs{"Input": input})
	if err != nil {
		return nil, fmt.Errorf("generate list error: %w", err)
	}

	schema := map[string]any{
		"type": "array",
		"items": map[string]any{
			"type":        "string",
			"description": "single item",
		},
	}

	response, err := assistant.StructuredAsk(ctx, "You build lists of items from text inputs", *prompt, schema)
	if err != nil {
		return nil, fmt.Errorf("generate list error: assistant structured ask: %w", err)
	}

	var items []string

	if err := json.Unmarshal(response, &items); err != nil {
		return nil, fmt.Errorf("generate list error: unmarshal json: %w", err)
	}

	return items, nil
}

func DocumentLength(doc *models.Document) int {
	length := 0

	for _, section := range doc.Sections {
		for _, paragraph := range section.Paragraphs {
			length += len(paragraph)
		}
	}

	return length
}
