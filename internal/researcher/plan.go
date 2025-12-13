package researcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

const (
	PlanSystemPrompt = `
		You are an expert Research Planner. Your sole task is
		to take a research description and break the topic down
		into a list of subtopics for a research report.
		`

	PlanPrompt = `
		# Research Topic
		{{.Topic}}

		# Goal
		Provide a list of 4 to 7 short subtopics for the give topic.

		# Tasks
		1. Perform an initial web search on the provided topis to gather some context.
		2. Decide on a series of 4 to 7 subtopics based on the description and initial research.
		3. For each subtopic, provide:
			- a title for the subtopic
			- a short description of the subtopic
		Present the result in a clear, readable text format. Do not use any HTML, markdown, or JSON.
		`
)

func (p *Pipeline) Plan(ctx context.Context, topic string, out chan<- Section) error {
	defer close(out)

	prompt, err := BuildPrompt(PlanPrompt, PromptArgs{
		"Topic": topic,
	})
	if err != nil {
		return fmt.Errorf("create plan error: %w", err)
	}

	response, err := p.assistant.Ask(ctx, PlanSystemPrompt, *prompt)
	if err != nil {
		return fmt.Errorf("create plan error: assistant ask: %w", err)
	}

	schema := map[string]any{
		"type":        "array",
		"description": "list of subtopics",
		"items": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title": map[string]any{
					"type":        "string",
					"description": "title of the subtopic",
				},
				"summary": map[string]any{
					"type":        "string",
					"description": "summary of the subtopic",
				},
			},
			"required": []string{"title", "summary"},
		},
	}

	structuredPrompt := "Extract the list of subtopics, title and summary, from the following text.\n" + *response

	responseJson, err := p.assistant.StructuredAsk(ctx, PlanSystemPrompt, structuredPrompt, schema)
	if err != nil {
		return fmt.Errorf("generate plan error: assistant structured ask: %w", err)
	}

	var sections []Section

	if err := json.Unmarshal(responseJson, &sections); err != nil {
		return fmt.Errorf("generate plan error: unmarshal json: %w", err)
	}

	slog.Info("plan",
		slog.String("topic", topic),
		slog.Int("sections", len(sections)),
	)

	for _, section := range sections {
		section.Topic = topic

		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- section:
		}
	}

	return nil
}
