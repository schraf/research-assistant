package researcher

import (
	"context"
	"log/slog"
	"strings"

	"github.com/schraf/assistant/pkg/models"
)

func Structure(ctx context.Context, report string) (models.Document, error) {
	var doc models.Document

	lines := strings.Split(report, "\n")
	var currentSectionTitle string
	var currentSectionBody strings.Builder

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "# ") && doc.Title == "" {
			doc.Title = strings.TrimPrefix(trimmedLine, "# ")
			continue
		}

		if strings.HasPrefix(trimmedLine, "## ") {
			if currentSectionTitle != "" {
				doc.AddSection(currentSectionTitle, strings.TrimSpace(currentSectionBody.String()))
			} else if doc.Title == "" && currentSectionBody.Len() > 0 {
				// Edge case: if there's text before the first section or title
			}
			currentSectionTitle = strings.TrimPrefix(trimmedLine, "## ")
			currentSectionBody.Reset()
			continue
		}

		currentSectionBody.WriteString(line + "\n")
	}

	if currentSectionTitle != "" {
		doc.AddSection(currentSectionTitle, strings.TrimSpace(currentSectionBody.String()))
	}

	doc.Clean()

	slog.Info("structured_document",
		slog.String("title", doc.Title),
		slog.Int("sections", len(doc.Sections)),
	)

	return doc, nil
}
