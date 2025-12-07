package researcher

import (
	"context"
	"log/slog"
)

const (
	AnalyzeKnowledgeSystemPrompt = `
		You are an expert Research Analyist. Your sole task is
		to review the information that other researchers have
		gathered.
		`

	AnalyzeKnowledgePrompt = `
		## Topic
		{{.Topic}}

		## Researched Information
		{{.Knowledge}}
		
		## Goal 
		Review the information gathered and see if there are any gaps in
		information about the topic that would require further information
		before we can synthesize a report about this topic. If you find
		that there are gaps in information, perform a series of web searches
		and return only new findings that we should add to our research.
		`
)

func (p *Pipeline) AnalyzeKnowledge(ctx context.Context, topic string, in <-chan string, out chan<- string) error {
	defer close(out)

	for knowledge := range in {
		slog.Info("analyzing_knowledge",
			slog.String("topic", topic),
			slog.String("knowledge", knowledge),
		)

		prompt, err := BuildPrompt(AnalyzeKnowledgePrompt, PromptArgs{
			"Topic":     topic,
			"Knowledge": knowledge,
		})
		if err != nil {
			return err
		}

		analysis, err := p.assistant.Ask(ctx, AnalyzeKnowledgeSystemPrompt, *prompt)
		if err != nil {
			return err
		}

		slog.Info("analyzed_knowledge",
			slog.String("topic", topic),
			slog.String("analysis", *analysis),
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- (knowledge + "\n" + *analysis):
		}
	}

	return nil
}
