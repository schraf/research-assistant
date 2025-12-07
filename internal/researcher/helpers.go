package researcher

import (
	"context"
	"encoding/json"

	"github.com/schraf/assistant/pkg/models"
)

const (
	ListPrompt = `
		Given the following list. Return a structured list that follows the provided schema.
		{{.Input}}
		`
)

func GenerateList(ctx context.Context, assistant models.Assistant, input string) ([]string, error) {
	prompt, err := BuildPrompt(ListPrompt, PromptArgs{"Input": input})
	if err != nil {
		return nil, err
	}

	schema := map[string]any{
		"type": "array",
		"items": map[string]any{
			"type": "string",
		},
	}

	response, err := assistant.StructuredAsk(ctx, "You build lists of items from text inputs", *prompt, schema)
	if err != nil {
		return nil, err
	}

	sanitized := Sanitize(response)

	var items []string

	if err := json.Unmarshal(sanitized, &items); err != nil {
		return nil, err
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

func Sanitize(input json.RawMessage) json.RawMessage {
	if json.Valid(input) {
		return input
	}

	validEscapeChars := map[byte]bool{
		'"':  true, // \"
		'\\': true, // \\
		'/':  true, // \/
		'b':  true, // \b
		'f':  true, // \f
		'n':  true, // \n
		'r':  true, // \r
		't':  true, // \t
	}

	result := make([]byte, 0, len(input)*2) // Pre-allocate with some extra capacity

	for i := 0; i < len(input); i++ {
		if input[i] == '\\' && i+1 < len(input) {
			nextChar := input[i+1]

			if validEscapeChars[nextChar] {
				result = append(result, input[i], input[i+1])
				i++ // Skip the next character
				continue
			}

			result = append(result, '\\', '\\', nextChar)
			i++ // Skip the next character
			continue
		}

		result = append(result, input[i])
	}

	return result
}
