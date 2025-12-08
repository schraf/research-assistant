package researcher

import (
	"context"

	"github.com/schraf/assistant/pkg/models"
	"golang.org/x/sync/errgroup"
)

func (p *Pipeline) CreateDocumentSection(ctx context.Context, in <-chan string, out chan<- models.DocumentSection, concurrency int) error {
	defer close(out)

	group, ctx := errgroup.WithContext(ctx)

	for i := 0; i < concurrency; i++ {
		group.Go(func() error {
			for subtopic := range in {
				topic := subtopic
				topicGroup, topicCtx := errgroup.WithContext(ctx)

				questions := make(chan string)
				information := make(chan string)
				knowledge := make(chan string)
				analysis := make(chan string)
				synthesis := make(chan models.DocumentSection, 1)

				topicGroup.Go(func() error {
					return p.GenerateQuestions(topicCtx, topic, questions)
				})

				topicGroup.Go(func() error {
					return p.ResearchQuestion(topicCtx, questions, information, 10)
				})

				topicGroup.Go(func() error {
					return p.BuildKnowledge(topicCtx, information, knowledge)
				})

				topicGroup.Go(func() error {
					return p.AnalyzeKnowledge(topicCtx, topic, knowledge, analysis)
				})

				topicGroup.Go(func() error {
					return p.SynthesizeKnowledge(topicCtx, topic, analysis, synthesis)
				})

				if err := topicGroup.Wait(); err != nil {
					return err
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				case out <- <-synthesis:
				}
			}

			return nil
		})
	}

	return group.Wait()
}
