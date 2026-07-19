package researcher

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
)

const (
	ComposeSystemPrompt = `
		You are an expert editor. Your role is to review content and 
		ensure that it takes a neutral stance and is clearly written.
		Remove any Markdown (except for # and ## headers), LaTeX, HTML tags, or any escape characters.
		`
)

func Compose(ctx context.Context, sections []Section) (string, error) {
	sortedSections := slices.Clone(sections)
	slices.SortFunc(sortedSections, func(a, b Section) int {
		return cmp.Compare(a.Index, b.Index)
	})

	// ╭────────────────────────────────────────────────────────────────────╮
	// │ Create a title for the document                                    │
	// ╰────────────────────────────────────────────────────────────────────╯

	var title string
	var summary string

	for _, section := range sortedSections {
		summary += "\n\n#" + section.Title + "\n\n" + section.Summary
	}

	schema := map[string]any{
		"type":        "string",
		"description": "document title",
	}

	responseJson, err := structuredAsk(ctx, "Create a document title", summary, schema)
	if err != nil {
		return "", fmt.Errorf("compose title: assistant ask: %w", err)
	}

	if err := json.Unmarshal(responseJson, &title); err != nil {
		return "", fmt.Errorf("compose title: unmarshal json: %w", err)
	}

	slog.Info("titled_document",
		slog.String("title", title),
	)

	// ╭────────────────────────────────────────────────────────────────────╮
	// │ Compose full document text                                         │
	// ╰────────────────────────────────────────────────────────────────────╯

	ctx = withDeepModel(ctx)

	content := "# " + title

	for _, section := range sortedSections {
		content += "\n\n## " + section.Title + "\n\n" + section.Body
	}

	report, err := ask(ctx, ComposeSystemPrompt, content)
	if err != nil {
		return "", fmt.Errorf("compose document: assistant ask: %w", err)
	}

	slog.Info("composed_document",
		slog.Int("length", len(*report)),
	)

	return *report, nil
}
