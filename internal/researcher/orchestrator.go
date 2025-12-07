package researcher

import (
	"context"

	"github.com/schraf/assistant/pkg/models"
)

func ResearchTopic(ctx context.Context, assistant models.Assistant, topic string, depth ResearchDepth) (*models.Document, error) {
	pipeline := NewPipeline(assistant)

	_, err := pipeline.Exec(ctx, topic)
	if err != nil {
		return nil, err
	}

	return &models.Document{
		Title:    "test",
		Author:   "test",
		Sections: []models.DocumentSection{},
	}, nil
}
