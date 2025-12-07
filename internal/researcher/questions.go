package researcher

import (
	"context"
	"log/slog"
)

const (
	QuestionsSystemPrompt = `
		You are an expert Research Planner. Your sole task is
		to take a research topic and search for what questions
		would need to be answered to fully understand the topic.
		`

	QuestionsPrompt = `
		# Research Topic
		{{.Topic}}

		# Goal
		Provide a list of research questions for the given topic.

		# Task
		1. Perform a series of web searches on the provided topic for this research.
		2. Build a list of questions that will need to be answered before writing on this topic.
		`
)

func (p *Pipeline) GenerateQuestions(ctx context.Context, topic string, out chan<- string) error {
	defer close(out)

	slog.Info("generate_questions",
		slog.String("topic", topic),
	)

	prompt, err := BuildPrompt(QuestionsPrompt, PromptArgs{
		"Topic": topic,
	})
	if err != nil {
		return err
	}

	response, err := p.assistant.Ask(ctx, QuestionsSystemPrompt, *prompt)
	if err != nil {
		return err
	}

	questions, err := GenerateList(ctx, p.assistant, *response)
	if err != nil {
		return err
	}

	slog.Info("generated_questions",
		slog.String("topic", topic),
		slog.Any("questions", questions),
	)

	for _, question := range questions {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- question:
		}
	}

	return nil
}
