package researcher

import (
	"context"
	"log/slog"
	"sync"

	"github.com/schraf/assistant/pkg/models"
)

type ResearchResult struct {
	Topic     string
	Knowledge []Knowledge
	Error     error
}

func ResearchTopic(ctx context.Context, assistant models.Assistant, topic string, depth ResearchDepth) (*models.Document, error) {
	plan, err := GenerateResearchPlan(ctx, assistant, topic, depth)
	if err != nil {
		return nil, err
	}

	resultsChan := make(chan ResearchResult, len(plan.ResearchItems))

	var group sync.WaitGroup

	for _, item := range plan.ResearchItems {
		group.Add(1)

		go func() {
			defer group.Done()

			slog.Info("starting_research",
				slog.String("topic", item.SubTopic),
			)

			resultsChan <- ResearchSubTopic(ctx, assistant, plan.Goal, item.SubTopic, item.Questions, depth)

			slog.Info("finished_research",
				slog.String("topic", item.SubTopic),
			)

		}()
	}

	group.Wait()
	close(resultsChan)

	results := []ResearchResult{}

	for result := range resultsChan {
		results = append(results, result)
	}

	doc, err := SynthesizeReport(ctx, assistant, topic, results, depth)
	if err != nil {
		return nil, err
	}

	doc, err = EditReport(ctx, assistant, doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func ResearchSubTopic(ctx context.Context, assistant models.Assistant, goal string, topic string, questions []string, depth ResearchDepth) ResearchResult {
	result := ResearchResult{
		Topic:     topic,
		Knowledge: []Knowledge{},
	}

	maxIterations := 0

	switch depth {
	case ResearchDepthShort:
		maxIterations = 0
	case ResearchDepthMedium:
		maxIterations = 2
	case ResearchDepthLong:
		maxIterations = 5
	}

	// build initial knowledge
	initialKnowledge, err := GenerateKnowledge(ctx, assistant, topic, questions, depth)
	if err != nil {
		result.Error = err
		return result
	}

	result.Knowledge = append(result.Knowledge, initialKnowledge...)

	for iteration := 0; iteration < maxIterations; iteration++ {
		questions, err := AnalyzeKnowledge(ctx, assistant, goal, topic, questions, result.Knowledge, depth)
		if err != nil {
			result.Error = err
			break
		}

		if len(questions) == 0 {
			break
		}

		newKnowledge, err := GenerateKnowledge(ctx, assistant, topic, questions, depth)
		if err != nil {
			result.Error = err
			break
		}

		result.Knowledge = append(result.Knowledge, newKnowledge...)
	}

	return result
}
