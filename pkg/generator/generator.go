package generator

import (
	"context"
	"fmt"

	"github.com/schraf/assistant/pkg/generators"
	"github.com/schraf/assistant/pkg/models"
	"github.com/schraf/research-assistant/internal/researcher"
)

func init() {
	generators.MustRegister("researcher", factory)
}

func factory(generators.Config) (models.ContentGenerator, error) {
	return &generator{}, nil
}

type generator struct{}

func (g *generator) Generate(ctx context.Context, request models.ContentRequest, assistant models.Assistant) (*models.Document, error) {
	topic, ok := request.Body["topic"].(string)
	if !ok {
		return nil, fmt.Errorf("no research topic")
	}

	depth, ok := request.Body["research_depth"].(researcher.ResearchDepth)
	if !ok {
		return nil, fmt.Errorf("no research depth")
	}

	if !depth.Validate() {
		return nil, fmt.Errorf("invalid research depth")
	}

	return researcher.ResearchTopic(ctx, assistant, topic, depth)
}
