package researcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/schraf/gemini-email/internal/models"
)

const (
	KnowledgeSystemPrompt = `
		You are an expert Researcher. Your sole task is
		to search the web to gather information about a
		topic by answering a set of questions.
		`

	KnowledgePrompt = `
		## Research Topic
		{{.ResearchTopic}}

		## Questions To Answer
		{{range $index, $question := .Questions}}
		{{$index}}. {{$question}}
		{{end}}
		
		## Goal 
		Provide information from searching the web that answers
		the questions. Be detailed in the information you provide
		and clearly specify what information is answering each 
		question.
		`

	KnowledgeStructureSystemPrompt = `
		You are an expert Researcher organizer. Your sole task is
		take information gathered and structure it into an
		organized list of pairs of topics and information.
		`

	KnowledgeStructurePrompt = `
		## Information Gathered
		{{.Information}}
		`
)

type Knowledge struct {
	Topic       string `json:"topic"`
	Information string `json:"information"`
}

func GenerateKnowledge(ctx context.Context, resources models.Resources, topic string, questions []string) ([]Knowledge, error) {
	slog.InfoContext(ctx, "generating_knowledge",
		slog.String("subtopic", topic),
	)

	prompt, err := BuildPrompt(KnowledgePrompt, PromptArgs{
		"ResearchTopic": topic,
		"Questions":     questions,
	})
	if err != nil {
		return nil, fmt.Errorf("failed building knowledge prompt: %w", err)
	}

	response, err := resources.Ask(ctx, KnowledgeSystemPrompt, *prompt)
	if err != nil {
		return nil, fmt.Errorf("failed gathering knowledge: %w", err)
	}

	structuredPrompt, err := BuildPrompt(KnowledgeStructurePrompt, PromptArgs{
		"Information": *response,
	})
	if err != nil {
		return nil, fmt.Errorf("failed building structured knowledge prompt: %w", err)
	}

	structuredResponse, err := resources.StructuredAsk(ctx, KnowledgeStructureSystemPrompt, *structuredPrompt, KnowledgeSchema())
	if err != nil {
		return nil, fmt.Errorf("failed structuring knowledge: %w", err)
	}

	var knowledge []Knowledge

	if err := json.Unmarshal(structuredResponse, &knowledge); err != nil {
		return nil, fmt.Errorf("failed parsing structured knowledge: %w", err)
	}

	return knowledge, nil
}

func KnowledgeSchema() map[string]any {
	return map[string]any{
		"type": "array",
		"items": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"topic": map[string]any{
					"type":        "string",
					"description": "A short description of the topic this information covers.",
				},
				"information": map[string]any{
					"type":        "string",
					"description": "A detailed report of information about this topic based on previously researched resources.",
				},
			},
			"required": []string{
				"topic",
				"information",
			},
		},
	}
}
