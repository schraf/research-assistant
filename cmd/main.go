package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/schraf/assistant/pkg/eval"
	"github.com/schraf/assistant/pkg/generators"
	"github.com/schraf/assistant/pkg/models"
	_ "github.com/schraf/research-assistant/pkg/generator"
)

func main() {
	topic := flag.String("topic", "", "Research topic (required)")
	flag.Parse()

	if *topic == "" {
		fmt.Fprintf(os.Stderr, "Error: topic is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Create request object
	request := models.ContentRequest{
		Body: map[string]any{
			"topic": *topic,
		},
	}

	ctx := context.Background()

	generator, err := generators.Create("researcher", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := eval.Evaluate(ctx, generator, request, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
