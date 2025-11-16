package utils

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// SetupLogger configures the default slog logger to output JSON logs
// to both the console and a log file with the specified log level.
func SetupLogger(logFile string, level slog.Level) error {
	// Create the log directory if it doesn't exist
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Open the log file for writing
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Create a multi-writer that writes to both console and file
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Create JSON handler for the multi-writer
	jsonHandler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: level,
	})

	// Set the default logger
	slog.SetDefault(slog.New(jsonHandler))

	return nil
}
