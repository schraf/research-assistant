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
	EditorSystemPrompt = `
		You are an expert editor. Your role reviews, revises, and optimizes 
		written content to ensure clarity, accuracy, and consistent formatting. 
		You improve readability while preserving the author's original voice.
		Ensure that later sections do no repeat concepts that were already
		covered in previous sections, instead reword them assuming the content
		is read from top to bottom.
		Remove any Markdown (except for # and ## headers), LaTeX, HTML tags, 
		or any escape characters.
		`
)

func Edit(ctx context.Context, sections []Section) (string, error) {
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
		return "", fmt.Errorf("edit title: assistant ask: %w", err)
	}

	if err := json.Unmarshal(responseJson, &title); err != nil {
		return "", fmt.Errorf("edit title: unmarshal json: %w", err)
	}

	slog.Info("titled_document",
		slog.String("title", title),
	)

	// ╭────────────────────────────────────────────────────────────────────╮
	// │ Edit full document text                                            │
	// ╰────────────────────────────────────────────────────────────────────╯

	ctx = withDeepModel(ctx)

	content := "# " + title

	for _, section := range sortedSections {
		content += "\n\n## " + section.Title + "\n\n" + section.Body
	}

	report, err := ask(ctx, EditorSystemPrompt, content)
	if err != nil {
		return "", fmt.Errorf("edit document: assistant ask: %w", err)
	}

	slog.Info("edited_document",
		slog.Int("length", len(*report)),
	)

	return *report, nil
}
