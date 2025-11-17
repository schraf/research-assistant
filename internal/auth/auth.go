package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
)

const (
	// SecretEnvVar is the environment variable name for the auth secret
	SecretEnvVar = "AUTH_SECRET"
	// AuthMessagesEnvVar is the environment variable name for valid token messages
	// Comma-separated list of messages that can be used to generate valid tokens
	AuthMessagesEnvVar = "AUTH_TOKEN_MESSAGES"
)

// ValidateToken validates a bearer token string against the secret stored in
// the environment variable. The token should be in the format "Bearer <token>"
// or just the token itself. Returns true if the token is valid, false otherwise.
//
// The function checks the token against messages specified in the AUTH_TOKEN_MESSAGES
// environment variable (comma-separated list). This allows multiple tokens to be valid
// for the same secret by using different messages when generating tokens.
func ValidateToken(bearerToken string) bool {
	// Get the secret from environment variable
	secret := os.Getenv(SecretEnvVar)
	if secret == "" {
		return false
	}

	// Remove "Bearer " prefix if present
	token := strings.TrimSpace(bearerToken)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = strings.TrimSpace(token[7:])
	}

	if token == "" {
		return false
	}

	// Get list of messages to validate against
	messages := getAuthMessages()
	if len(messages) == 0 {
		return false
	}

	// Check token against all possible messages
	for _, message := range messages {
		expectedToken := generateToken(secret, message)
		if hmac.Equal([]byte(token), []byte(expectedToken)) {
			return true
		}
	}

	return false
}

// getAuthMessages returns a list of messages to validate tokens against.
// Messages are read from the AUTH_TOKEN_MESSAGES environment variable (comma-separated).
func getAuthMessages() []string {
	var messages []string

	// Get messages from environment variable
	authMessages := os.Getenv(AuthMessagesEnvVar)
	if authMessages != "" {
		// Split by comma and trim whitespace
		messageList := strings.Split(authMessages, ",")
		for _, msg := range messageList {
			trimmed := strings.TrimSpace(msg)
			if trimmed != "" {
				messages = append(messages, trimmed)
			}
		}
	}

	return messages
}

// GenerateToken generates a token using HMAC-SHA256 with the given secret and message.
func GenerateToken(secret, message string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// generateToken generates a token using HMAC-SHA256 with the given secret and message.
func generateToken(secret, message string) string {
	return GenerateToken(secret, message)
}
