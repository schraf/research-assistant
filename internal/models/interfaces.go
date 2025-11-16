package models

import (
	"context"
	"encoding/json"
)

type Resources interface {
	Ask(ctx context.Context, persona string, request string) (*string, error)
	StructuredAsk(ctx context.Context, persona string, request string, schema Schema) (json.RawMessage, error)
}
