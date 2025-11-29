package researcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/schraf/assistant/pkg/models"
)

const (
	AnalyzeKnowledgeSystemPrompt = `
		You are an expert Research Analyist. Your sole task is
		to review the information that other researchers have
		gathers and verify if it satisfies the questions for
		the given subtopic of the provided research goal.
		`

	AnalyzeKnowledgePrompt = `
		## Research Goal
		{{.ResearchGoal}}

		## Research Subtopic
		{{.ResearchTopic}}

		## Questions To Answer
		{{range $index, $question := .Questions}}
		{{$index}}. {{$question}}
		{{end}}

		## Information Gathered
		{{range $index, $knowledge := .Knowledge}}
		{{$index}}. {{$knowledge.Topic}}
		{{$knowledge.Information}}

		{{end}}
		
		## Goal 
		Review the information gathered and see if it sufficiently answers all
		of the questions in enough detail to fullfil the information gathering
		phase of research for this particular subtopic. If there are gaps in
		knowledge provide a list of further research questions for this
		subtopic.  
		`
)

func AnalyzeKnowledge(ctx context.Context, assistant models.Assistant, goal string, topic string, questions []string, knowledge []Knowledge, depth ResearchDepth) ([]string, error) {
	slog.InfoContext(ctx, "analyzing_knowledge",
		slog.String("topic", topic),
	)

	prompt, err := BuildPrompt(AnalyzeKnowledgePrompt, PromptArgs{
		"ResearchGoal":  goal,
		"ResearchTopic": topic,
		"Questions":     questions,
		"Knowledge":     knowledge,
	})
	if err != nil {
		return nil, fmt.Errorf("failed building analyze knowledge prompt: %w", err)
	}

	response, err := assistant.StructuredAsk(ctx, AnalyzeKnowledgeSystemPrompt, *prompt, AnalyzeKnowledgeSchema())
	if err != nil {
		return nil, fmt.Errorf("failed analyze knowledge request: %w", err)
	}

	var furtherQuestions []string

	if err := json.Unmarshal(response, &furtherQuestions); err != nil {
		return nil, fmt.Errorf("failed parsing analyze knowledge response: %w", err)
	}

	return furtherQuestions, nil
}

func AnalyzeKnowledgeSchema() map[string]any {
	return map[string]any{
		"type":        "array",
		"description": "A list of follow up research questions to be answered",
		"items": map[string]any{
			"type":        "string",
			"description": "Research question",
		},
	}
}
