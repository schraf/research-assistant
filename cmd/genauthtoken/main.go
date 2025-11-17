package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/schraf/research-assistant/internal/auth"
	"github.com/schraf/research-assistant/internal/utils"
)

func main() {
	if err := utils.LoadEnv(".env"); err != nil {
		slog.Error("load_env_failed",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	// Get default value from environment variable
	defaultSeed := os.Getenv(auth.SecretEnvVar)

	// Parse command-line flags
	seed := flag.String("seed", defaultSeed, "Secret seed for token generation (required, can also be set via AUTH_SECRET env var)")
	flag.Parse()

	// Validate required flag
	if *seed == "" {
		fmt.Fprintf(os.Stderr, "Error: -seed is required (or set %s environment variable)\n", auth.SecretEnvVar)
		flag.Usage()

		os.Exit(1)
	}

	// Get messages from environment variable
	authMessages := os.Getenv(auth.AuthMessagesEnvVar)
	if authMessages == "" {
		fmt.Fprintf(os.Stderr, "Error: %s environment variable is not set\n", auth.AuthMessagesEnvVar)
		os.Exit(1)
	}

	// Parse messages (comma-separated)
	messageList := strings.Split(authMessages, ",")
	var messages []string
	for _, msg := range messageList {
		trimmed := strings.TrimSpace(msg)
		if trimmed != "" {
			messages = append(messages, trimmed)
		}
	}

	if len(messages) == 0 {
		fmt.Fprintf(os.Stderr, "Error: %s contains no valid messages\n", auth.AuthMessagesEnvVar)
		os.Exit(1)
	}

	// Generate and output tokens for each message
	for _, message := range messages {
		token := auth.GenerateToken(*seed, message)
		fmt.Printf("%s,%s\n", message, token)
	}
}
