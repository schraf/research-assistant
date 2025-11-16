package gemini

import (
	"context"
	"encoding/json"

	"github.com/schraf/gemini-email/internal/models"
)

type Client interface {
	models.Resources

	EnableLogging(filename string) (closer func(), err error)
	SetSystemInstruction(instruction string)
	GenerateText(ctx context.Context, model ModelIdentifier, prompt string) (*string, error)
	GenerateJson(ctx context.Context, model ModelIdentifier, prompt string, schema map[string]any) (json.RawMessage, error)
	Chat(ctx context.Context, model ModelIdentifier, history ChatHistory, message string) (*string, error)
}
