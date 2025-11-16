package researcher

import (
	"context"
	"log/slog"
	"sync"

	"github.com/schraf/gemini-email/internal/models"
)

type ResearchResult struct {
	Topic     string
	Knowledge []Knowledge
	Error     error
}

func ResearchTopic(ctx context.Context, resources models.Resources, topic string) (*ResearchReport, error) {
	plan, err := GenerateResearchPlan(ctx, resources, topic)
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

			resultsChan <- ResearchSubTopic(ctx, resources, plan.Goal, item.SubTopic, item.Questions, 5)

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

	report, err := SynthesizeReport(ctx, resources, topic, results)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func ResearchSubTopic(ctx context.Context, resources models.Resources, goal string, topic string, questions []string, maxIterations int) ResearchResult {
	result := ResearchResult{
		Topic:     topic,
		Knowledge: []Knowledge{},
	}

	for iteration := 0; iteration < maxIterations; iteration++ {
		newKnowledge, err := GenerateKnowledge(ctx, resources, topic, questions)
		if err != nil {
			result.Error = err
			break
		}

		result.Knowledge = append(result.Knowledge, newKnowledge...)

		moreQuestions, err := AnalyzeKnowledge(ctx, resources, goal, topic, questions, result.Knowledge)
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
