package researcher

import (
	"context"
	"fmt"

	"github.com/schraf/assistant/pkg/models"
	"github.com/schraf/pipeline/v4"
	"github.com/schraf/pipeline/v4/stages"
)

func CreateDocument(ctx context.Context, assistant models.Assistant, topic string) (*models.Document, error) {
	ctx = withAssistant(ctx, assistant)

	// ╭────────────────────────────────────────────────────────────────────╮
	// │ Pipeline Stages                                                    │
	// ╰────────────────────────────────────────────────────────────────────╯

	planStage := stages.ExpandStage[string, Section]{
		Name:     "plan",
		Expander: Plan,
	}

	researchStage := stages.ParallelTransformStage[Section, Section]{
		Name:        "research_section",
		Workers:     6,
		Transformer: Research,
	}

	synthesisStage := stages.ParallelTransformStage[Section, Section]{
		Name:        "synthesize_section",
		Workers:     6,
		Transformer: Synthesize,
	}

	aggregateStage := stages.AggregateStage[Section]{
		Name: "aggregate_sections",
	}

	composeStage := stages.TransformStage[[]Section, string]{
		Name:        "compose_document",
		Transformer: Compose,
	}

	structureStage := stages.TransformStage[string, models.Document]{
		Name:        "structure_document",
		Transformer: Structure,
	}

	// ╭────────────────────────────────────────────────────────────────────╮
	// │ Compose Pipeline                                                   │
	// ╰────────────────────────────────────────────────────────────────────╯

	cfg := pipeline.Config[string, models.Document]{
		Name:             "researcher",
		OutputBufferSize: 1,
		Composer: func(c pipeline.Composer[string, models.Document]) error {
			stage1 := planStage.Create(c.Context(), c.Inputs().At(0))
			stage2 := researchStage.Create(c.Context(), stage1)
			stage3 := synthesisStage.Create(c.Context(), stage2)
			stage4 := aggregateStage.Create(c.Context(), stage3)
			stage5 := composeStage.Create(c.Context(), stage4)
			stage6 := structureStage.Create(c.Context(), stage5)

			out := stage6

			return c.Outputs().Link(c.Context(), 0, out)
		},
	}

	pipe, pipeCtx, err := pipeline.NewPipeline(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// ╭────────────────────────────────────────────────────────────────────╮
	// │ Execute Pipeline                                                   │
	// ╰────────────────────────────────────────────────────────────────────╯

	pipe.Inputs().Send(pipeCtx, 0, topic)
	pipe.CloseAllInputs()

	if err := pipe.Wait(); err != nil {
		return nil, err
	}

	documents := pipe.Outputs().SinkAt(ctx, 0)

	if len(documents) == 0 {
		return nil, fmt.Errorf("pipeline result empty")
	}

	return &documents[0], nil
}
