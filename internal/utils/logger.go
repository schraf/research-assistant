package utils

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// cloudLoggingHandler wraps a JSON handler and converts slog levels to
// Google Cloud Logging severity levels for proper log level detection.
type cloudLoggingHandler struct {
	slog.Handler
}

// Handle converts the log record to use "severity" instead of "level"
// and maps slog levels to Cloud Logging severity values.
func (h *cloudLoggingHandler) Handle(ctx context.Context, r slog.Record) error {
	// Create a new record with the severity field
	attrs := make([]slog.Attr, 0, r.NumAttrs()+1)
	
	// Map slog level to Cloud Logging severity
	severity := levelToSeverity(r.Level)
	attrs = append(attrs, slog.String("severity", severity))
	
	// Copy all other attributes, but skip the default "level" attribute
	r.Attrs(func(a slog.Attr) bool {
		if a.Key != "level" {
			attrs = append(attrs, a)
		}
		return true
	})
	
	// Create a new record with the modified attributes
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	for _, attr := range attrs {
		newRecord.AddAttrs(attr)
	}
	
	return h.Handler.Handle(ctx, newRecord)
}

// levelToSeverity maps slog.Level to Google Cloud Logging severity strings.
func levelToSeverity(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "ERROR"
	case level >= slog.LevelWarn:
		return "WARNING"
	case level >= slog.LevelInfo:
		return "INFO"
	case level >= slog.LevelDebug:
		return "DEBUG"
	default:
		return "DEFAULT"
	}
}

// isCloudRun detects if the application is running in Google Cloud Run.
func isCloudRun() bool {
	// Cloud Run sets K_SERVICE environment variable
	return os.Getenv("K_SERVICE") != ""
}

// SetupLogger configures the default slog logger to output JSON logs
// with Google Cloud Logging compatible severity levels.
// When running in Cloud Run, logs only go to stdout.
// Otherwise, logs go to both console and a log file.
func SetupLogger(logFile string, level slog.Level) error {
	var writer io.Writer
	
	if isCloudRun() {
		// In Cloud Run, only write to stdout (Cloud Logging captures stdout/stderr)
		writer = os.Stdout
	} else {
		// Create the log directory if it doesn't exist
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		// Open the log file for writing
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		// Create a multi-writer that writes to both console and file
		writer = io.MultiWriter(os.Stdout, file)
	}

	// Create JSON handler
	jsonHandler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: level,
		// ReplaceAttr removes the default "level" field since we'll add "severity" instead
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove the default "level" attribute - we'll add "severity" in the handler
			if a.Key == "level" {
				return slog.Attr{}
			}
			return a
		},
	})

	// Wrap with Cloud Logging compatible handler
	cloudHandler := &cloudLoggingHandler{Handler: jsonHandler}

	// Set the default logger
	slog.SetDefault(slog.New(cloudHandler))

	return nil
}
