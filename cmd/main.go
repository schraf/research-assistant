package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/schraf/assistant/pkg/eval"
	"github.com/schraf/assistant/pkg/generators"
	"github.com/schraf/assistant/pkg/models"
	"github.com/schraf/research-assistant/internal/researcher"
	_ "github.com/schraf/research-assistant/pkg/generator"
)

func main() {
	topic := flag.String("topic", "", "Research topic (required)")
	depthString := flag.String("depth", "basic", "Research depth: basic, medium, or long (default: basic)")
	model := flag.String("model", "", "Model to use for evaluation")
	flag.Parse()

	if *topic == "" {
		fmt.Fprintf(os.Stderr, "Error: topic is required\n")
		flag.Usage()
		os.Exit(1)
	}

	var depth researcher.ResearchDepth
	switch *depthString {
	case "basic":
		depth = researcher.ResearchDepthShort
	case "medium":
		depth = researcher.ResearchDepthMedium
	case "long":
		depth = researcher.ResearchDepthLong
	default:
		fmt.Fprintf(os.Stderr, "Error: invalid depth '%s'. Must be one of: basic, medium, long\n", *depthString)
		os.Exit(1)
	}

	// Create request object
	request := models.ContentRequest{
		Body: map[string]any{
			"topic":          *topic,
			"research_depth": depth,
		},
	}

	ctx := context.Background()

	generator, err := generators.Create("researcher", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := eval.Evaluate(ctx, generator, request, *model); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
