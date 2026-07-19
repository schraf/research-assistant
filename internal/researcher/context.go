package researcher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/schraf/assistant/pkg/models"
)

type contextKey int

const (
	assistantContextKey contextKey = iota
	modelContextKey
)

func withAssistant(ctx context.Context, assistant models.Assistant) context.Context {
	return context.WithValue(ctx, assistantContextKey, assistant)
}

func withLiteModel(ctx context.Context) context.Context {
	return context.WithValue(ctx, modelContextKey, models.ModelLite)
}

func withDeepModel(ctx context.Context) context.Context {
	return context.WithValue(ctx, modelContextKey, models.ModelDeep)
}

func ask(ctx context.Context, persona string, request string) (*string, error) {
	assistant, ok := ctx.Value(assistantContextKey).(models.Assistant)
	if !ok {
		return nil, fmt.Errorf("no assistant in context")
	}

	model, ok := ctx.Value(modelContextKey).(models.ModelType)
	if ok {
		ctx = assistant.WithModel(ctx, model)
	}

	return assistant.Ask(ctx, persona, request)
}

func structuredAsk(ctx context.Context, persona string, request string, schema map[string]any) (json.RawMessage, error) {
	assistant, ok := ctx.Value(assistantContextKey).(models.Assistant)
	if !ok {
		return nil, fmt.Errorf("no assistant in context")
	}

	model, ok := ctx.Value(modelContextKey).(models.ModelType)
	if ok {
		ctx = assistant.WithModel(ctx, model)
	}

	return assistant.StructuredAsk(ctx, persona, request, schema)
}
