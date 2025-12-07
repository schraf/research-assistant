package researcher

import (
	"context"
	"log/slog"
)

const (
	KnowledgeSystemPrompt = `
		You are an expert Researcher. Your sole task is
		to read through raw researched information and organize
		it into a comprehensive and conheret knowledge base.
		`

	KnowledgePrompt = `
		## Information
		{{range $index, $information := .Information}}
		{{$index}}. {{$information}}
		{{end}}
		
		## Goal 
		Analyze the raw information provided and produce a well
		organized knowledge base. The knowledge base should be
		a series of paragraphs. There is no need to add headings.
		`
)

func (p *Pipeline) BuildKnowledge(ctx context.Context, in <-chan string, out chan<- string) error {
	defer close(out)

	aggregated := []string{}

	for information := range in {
		aggregated = append(aggregated, information)
	}

	slog.Info("building_knowledge",
		slog.Any("information", aggregated),
	)

	prompt, err := BuildPrompt(KnowledgePrompt, PromptArgs{
		"Information": aggregated,
	})
	if err != nil {
		return err
	}

	knowledge, err := p.assistant.Ask(ctx, KnowledgeSystemPrompt, *prompt)
	if err != nil {
		return err
	}

	slog.Info("built_knowledge",
		slog.Any("information", aggregated),
		slog.String("knowledge", *knowledge),
	)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case out <- *knowledge:
	}

	return nil
}
