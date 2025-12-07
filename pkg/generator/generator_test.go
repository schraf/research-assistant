package generator

import (
	"context"
	"os"
	"testing"

	"github.com/schraf/assistant/pkg/eval"
	"github.com/schraf/assistant/pkg/generators"
	"github.com/schraf/assistant/pkg/models"
	"github.com/schraf/research-assistant/internal/researcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	os.Setenv("ASSISTANT_PROVIDER", "mock")

	request := models.ContentRequest{
		Body: map[string]any{
			"topic":          "test",
			"research_depth": researcher.ResearchDepthShort,
		},
	}

	ctx := context.Background()

	generator, err := generators.Create("researcher", nil)
	require.NoError(t, err)

	err = eval.Evaluate(ctx, generator, request, "mock-model")
	assert.NoError(t, err)
}
