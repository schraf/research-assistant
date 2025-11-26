package researcher

import (
	"context"
	"log/slog"
	"sync"

	"github.com/schraf/research-assistant/internal/models"
)

type ResearchResult struct {
	Topic     string
	Knowledge []Knowledge
	Error     error
}

func ResearchTopic(ctx context.Context, logger *slog.Logger, resources models.Resources, topic string, mode models.ResourceMode, depth models.ResearchDepth) (*ResearchReport, error) {
	plan, err := GenerateResearchPlan(ctx, logger, resources, topic, mode, depth)
	if err != nil {
		return nil, err
	}

	resultsChan := make(chan ResearchResult, len(plan.ResearchItems))

	var group sync.WaitGroup

	for _, item := range plan.ResearchItems {
		group.Add(1)

		go func() {
			defer group.Done()

			logger.Info("starting_research",
				slog.String("topic", item.SubTopic),
			)

			resultsChan <- ResearchSubTopic(ctx, logger, resources, plan.Goal, item.SubTopic, item.Questions, mode, depth)

			logger.Info("finished_research",
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

	report, err := SynthesizeReport(ctx, logger, resources, topic, results, mode, depth)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func ResearchSubTopic(ctx context.Context, logger *slog.Logger, resources models.Resources, goal string, topic string, questions []string, mode models.ResourceMode, depth models.ResearchDepth) ResearchResult {
	result := ResearchResult{
		Topic:     topic,
		Knowledge: []Knowledge{},
	}

	maxIterations := 0

	switch depth {
	case models.ResearchDepthShort:
		maxIterations = 0
	case models.ResearchDepthMedium:
		maxIterations = 2
	case models.ResearchDepthLong:
		maxIterations = 5
	}

	// build initial knowledge
	initialKnowledge, err := GenerateKnowledge(ctx, logger, resources, topic, questions, mode, depth)
	if err != nil {
		result.Error = err
		return result
	}

	result.Knowledge = append(result.Knowledge, initialKnowledge...)

	for iteration := 0; iteration < maxIterations; iteration++ {
		questions, err := AnalyzeKnowledge(ctx, logger, resources, goal, topic, questions, result.Knowledge, mode, depth)
		if err != nil {
			result.Error = err
			break
		}

		if len(questions) == 0 {
			break
		}

		newKnowledge, err := GenerateKnowledge(ctx, logger, resources, topic, questions, mode, depth)
		if err != nil {
			result.Error = err
			break
		}

		result.Knowledge = append(result.Knowledge, newKnowledge...)
	}

	return result
}
