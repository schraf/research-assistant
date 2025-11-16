package auth

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2/google"
)

func NewClient(ctx context.Context, scopes []string) (*http.Client, error) {
	client, err := google.DefaultClient(ctx, scopes...)
	if err != nil {
		return nil, fmt.Errorf("unable to create auth client: %w", err)
	}

	return client, nil
}
