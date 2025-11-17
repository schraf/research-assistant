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

func ResearchTopic(ctx context.Context, logger *slog.Logger, resources models.Resources, topic string) (*ResearchReport, error) {
	plan, err := GenerateResearchPlan(ctx, logger, resources, topic)
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

			resultsChan <- ResearchSubTopic(ctx, logger, resources, plan.Goal, item.SubTopic, item.Questions, 5)

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

	report, err := SynthesizeReport(ctx, logger, resources, topic, results)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func ResearchSubTopic(ctx context.Context, logger *slog.Logger, resources models.Resources, goal string, topic string, questions []string, maxIterations int) ResearchResult {
	result := ResearchResult{
		Topic:     topic,
		Knowledge: []Knowledge{},
	}

	for iteration := 0; iteration < maxIterations; iteration++ {
		newKnowledge, err := GenerateKnowledge(ctx, logger, resources, topic, questions)
		if err != nil {
			result.Error = err
			break
		}

		result.Knowledge = append(result.Knowledge, newKnowledge...)

		moreQuestions, err := AnalyzeKnowledge(ctx, logger, resources, goal, topic, questions, result.Knowledge)
		if err != nil {
			result.Error = err
			break
		}

		if len(moreQuestions) == 0 {
			break
		}

		questions = moreQuestions
	}

	return result
}
