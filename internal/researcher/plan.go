package researcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/schraf/gemini-email/internal/models"
)

const (
	ResearchPlanSystemPrompt = `
		You are an expert Research Planner. Your sole task is
		to take a research description and convert it into a
		fully structured, executable research plan.
		`

	ResearchPlanPrompt = `
		## Research Description
		{{.ResearchTopic}}

		## Task Constraints
		
		### Goal 
		Create a complete plan to fully research the description.

		### Research Items
		The "research_items" array must list **all necessary individual
		sub-topics** required to answer the main goal. Each sub-topic **must be
		paired** with a list of the most effective initial questions that must
		be searched to find its answer.
		`
)

type ResearchItem struct {
	SubTopic  string   `json:"subtopic"`
	Questions []string `json:"questions"`
}

type ResearchPlan struct {
	Goal          string         `json:"goal"`
	ResearchItems []ResearchItem `json:"research_items"`
}

func GenerateResearchPlan(ctx context.Context, resources models.Resources, topic string) (*ResearchPlan, error) {
	slog.InfoContext(ctx, "generating_research_plan")

	prompt, err := BuildPrompt(ResearchPlanPrompt, PromptArgs{"ResearchTopic": topic})
	if err != nil {
		return nil, fmt.Errorf("failed building research prompt: %w", err)
	}

	response, err := resources.StructuredAsk(ctx, ResearchPlanSystemPrompt, *prompt, ResearchPlanSchema())
	if err != nil {
		return nil, fmt.Errorf("failed building research plan: %w", err)
	}

	var plan ResearchPlan

	if err := json.Unmarshal(response, &plan); err != nil {
		return nil, fmt.Errorf("failed parsing research plan: %w", err)
	}

	return &plan, nil
}

func ResearchPlanSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"goal": map[string]any{
				"type":        "string",
				"description": "The overall detailed research goals to evaluate when work is complete.",
			},
			"research_items": map[string]any{
				"type":        "array",
				"description": "A list of key topics and their questions that will need answers.",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"subtopic": map[string]any{
							"type":        "string",
							"description": "A research subtopic that must be answered to complete the research goals",
						},
						"questions": map[string]any{
							"type":        "array",
							"description": "A list of questions that need futher research for the subtopic",
							"items": map[string]any{
								"type": "string",
							},
						},
					},
					"required": []string{
						"subtopic",
						"questions",
					},
				},
			},
		},
		"required": []string{
			"goal",
			"research_items",
		},
	}
}
