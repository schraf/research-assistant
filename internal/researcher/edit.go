package researcher

import (
	"context"
	"fmt"
	"log/slog"
)

const (
	EditSystemPrompt = `
		You are an expert editor. Your role is to review content and 
		ensure that it takes a neutral stance and is clearly written.
		Remove any Markdown, LaTeX, HTML tags, or any escape characters. The
		content should not include any headings. 
		`
)

func Edit(ctx context.Context, section Section) (*Section, error) {
	body, err := ask(ctx, EditSystemPrompt, section.Body)
	if err != nil {
		return nil, fmt.Errorf("edit error: assistant ask: %w", err)
	}

	section.Body = *body

	slog.Info("edited_section",
		slog.String("section", section.Title),
		slog.Int("length", len(section.Body)),
	)

	return &section, nil
}
