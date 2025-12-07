package researcher

import (
	"context"
	"log/slog"
)

const (
	SubtopicsSystemPrompt = `
		You are an expert Research Planner. Your sole task is
		to take a research description and break the topic down
		into a list of subtopics.
		`

	SubtopicsPrompt = `
		# Research Topic
		{{.Topic}}

		# Goal
		Provide a list of subtopics for the give topic.

		# Tasks
		1. Perform an initial web search on the provided topis to gather some context.
		2. Decide on a series of subtopics based on the description and initial research.
		`
)

func (p *Pipeline) CreateSubtopics(ctx context.Context, topic string, out chan<- string) error {
	defer close(out)

	slog.Info("creating_subtopics",
		slog.String("topic", topic),
	)

	prompt, err := BuildPrompt(SubtopicsPrompt, PromptArgs{
		"Topic": topic,
	})
	if err != nil {
		return err
	}

	response, err := p.assistant.Ask(ctx, SubtopicsSystemPrompt, *prompt)
	if err != nil {
		return err
	}

	subtopics, err := GenerateList(ctx, p.assistant, *response)
	if err != nil {
		return err
	}

	slog.Info("created_subtopics",
		slog.String("topic", topic),
		slog.Any("subtopics", subtopics),
	)

	for _, subtopic := range subtopics {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- subtopic:
		}
	}

	return nil
}
