package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
)

// generateTestToken generates a token for testing purposes
func generateTestToken(secret, message string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestValidateToken(t *testing.T) {
	// Save and restore the original environment variables
	originalSecret := os.Getenv(SecretEnvVar)
	originalMessages := os.Getenv(AuthMessagesEnvVar)
	defer func() {
		if originalSecret != "" {
			os.Setenv(SecretEnvVar, originalSecret)
		} else {
			os.Unsetenv(SecretEnvVar)
		}
		if originalMessages != "" {
			os.Setenv(AuthMessagesEnvVar, originalMessages)
		} else {
			os.Unsetenv(AuthMessagesEnvVar)
		}
	}()

	tests := []struct {
		name        string
		secret      string
		messages    string
		token       string
		expectValid bool
	}{
		{
			name:        "valid token with Bearer prefix",
			secret:      "test-secret-123",
			messages:    "auth-token",
			token:       "Bearer " + generateTestToken("test-secret-123", "auth-token"),
			expectValid: true,
		},
		{
			name:        "valid token without Bearer prefix",
			secret:      "test-secret-123",
			messages:    "auth-token",
			token:       generateTestToken("test-secret-123", "auth-token"),
			expectValid: true,
		},
		{
			name:        "valid token with lowercase bearer prefix",
			secret:      "test-secret-123",
			messages:    "auth-token",
			token:       "bearer " + generateTestToken("test-secret-123", "auth-token"),
			expectValid: true,
		},
		{
			name:        "valid token with mixed case bearer prefix",
			secret:      "test-secret-123",
			messages:    "auth-token",
			token:       "BeArEr " + generateTestToken("test-secret-123", "auth-token"),
			expectValid: true,
		},
		{
			name:        "invalid token",
			secret:      "test-secret-123",
			messages:    "auth-token",
			token:       "invalid-token-12345",
			expectValid: false,
		},
		{
			name:        "wrong token for secret",
			secret:      "test-secret-123",
			messages:    "auth-token",
			token:       generateTestToken("different-secret", "auth-token"),
			expectValid: false,
		},
		{
			name:        "empty token",
			secret:      "test-secret-123",
			messages:    "auth-token",
			token:       "",
			expectValid: false,
		},
		{
			name:        "empty token with Bearer prefix",
			secret:      "test-secret-123",
			messages:    "auth-token",
			token:       "Bearer ",
			expectValid: false,
		},
		{
			name:        "missing AUTH_SECRET environment variable",
			secret:      "",
			messages:    "auth-token",
			token:       generateTestToken("test-secret-123", "auth-token"),
			expectValid: false,
		},
		{
			name:        "missing AUTH_TOKEN_MESSAGES environment variable",
			secret:      "test-secret-123",
			messages:    "",
			token:       generateTestToken("test-secret-123", "auth-token"),
			expectValid: false,
		},
		{
			name:        "token with extra whitespace",
			secret:      "test-secret-123",
			messages:    "auth-token",
			token:       "  Bearer  " + generateTestToken("test-secret-123", "auth-token") + "  ",
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			if tt.secret != "" {
				os.Setenv(SecretEnvVar, tt.secret)
			} else {
				os.Unsetenv(SecretEnvVar)
			}
			if tt.messages != "" {
				os.Setenv(AuthMessagesEnvVar, tt.messages)
			} else {
				os.Unsetenv(AuthMessagesEnvVar)
			}

			// Validate token
			result := ValidateToken(tt.token)

			// Check result
			if result != tt.expectValid {
				t.Errorf("ValidateToken() = %v, want %v", result, tt.expectValid)
			}
		})
	}
}

func TestValidateToken_EmptySecret(t *testing.T) {
	// Save and restore the original environment variables
	originalSecret := os.Getenv(SecretEnvVar)
	originalMessages := os.Getenv(AuthMessagesEnvVar)
	defer func() {
		if originalSecret != "" {
			os.Setenv(SecretEnvVar, originalSecret)
		} else {
			os.Unsetenv(SecretEnvVar)
		}
		if originalMessages != "" {
			os.Setenv(AuthMessagesEnvVar, originalMessages)
		} else {
			os.Unsetenv(AuthMessagesEnvVar)
		}
	}()

	// Unset the environment variables
	os.Unsetenv(SecretEnvVar)
	os.Setenv(AuthMessagesEnvVar, "auth-token")

	// Any token should be invalid when secret is missing
	validToken := generateTestToken("some-secret", "auth-token")
	if ValidateToken(validToken) {
		t.Error("ValidateToken() should return false when AUTH_SECRET is not set")
	}

	if ValidateToken("Bearer " + validToken) {
		t.Error("ValidateToken() should return false when AUTH_SECRET is not set")
	}
}

func TestValidateToken_DifferentSecrets(t *testing.T) {
	// Save and restore the original environment variables
	originalSecret := os.Getenv(SecretEnvVar)
	originalMessages := os.Getenv(AuthMessagesEnvVar)
	defer func() {
		if originalSecret != "" {
			os.Setenv(SecretEnvVar, originalSecret)
		} else {
			os.Unsetenv(SecretEnvVar)
		}
		if originalMessages != "" {
			os.Setenv(AuthMessagesEnvVar, originalMessages)
		} else {
			os.Unsetenv(AuthMessagesEnvVar)
		}
	}()

	secret1 := "secret-one"
	secret2 := "secret-two"
	message := "auth-token"

	// Generate tokens for each secret
	token1 := generateTestToken(secret1, message)
	token2 := generateTestToken(secret2, message)

	// Set secret1 and messages
	os.Setenv(SecretEnvVar, secret1)
	os.Setenv(AuthMessagesEnvVar, message)

	// Token1 should be valid
	if !ValidateToken(token1) {
		t.Error("ValidateToken() should return true for token generated with secret1")
	}

	// Token2 should be invalid
	if ValidateToken(token2) {
		t.Error("ValidateToken() should return false for token generated with different secret")
	}

	// Switch to secret2
	os.Setenv(SecretEnvVar, secret2)

	// Token1 should now be invalid
	if ValidateToken(token1) {
		t.Error("ValidateToken() should return false for token generated with different secret")
	}

	// Token2 should now be valid
	if !ValidateToken(token2) {
		t.Error("ValidateToken() should return true for token generated with secret2")
	}
}

func TestValidateToken_MultipleMessages(t *testing.T) {
	// Save and restore the original environment variables
	originalSecret := os.Getenv(SecretEnvVar)
	originalMessages := os.Getenv(AuthMessagesEnvVar)
	defer func() {
		if originalSecret != "" {
			os.Setenv(SecretEnvVar, originalSecret)
		} else {
			os.Unsetenv(SecretEnvVar)
		}
		if originalMessages != "" {
			os.Setenv(AuthMessagesEnvVar, originalMessages)
		} else {
			os.Unsetenv(AuthMessagesEnvVar)
		}
	}()

	secret := "test-secret-123"
	os.Setenv(SecretEnvVar, secret)

	// Test 1: Multiple messages should work
	os.Setenv(AuthMessagesEnvVar, "auth-token-1,auth-token-2,auth-token-3")
	
	token1 := generateTestToken(secret, "auth-token-1")
	token2 := generateTestToken(secret, "auth-token-2")
	token3 := generateTestToken(secret, "auth-token-3")

	if !ValidateToken(token1) {
		t.Error("ValidateToken() should return true for token1")
	}
	if !ValidateToken(token2) {
		t.Error("ValidateToken() should return true for token2")
	}
	if !ValidateToken(token3) {
		t.Error("ValidateToken() should return true for token3")
	}

	// Test 2: Invalid message should not work
	invalidToken := generateTestToken(secret, "invalid-message")
	if ValidateToken(invalidToken) {
		t.Error("ValidateToken() should return false for invalid message token")
	}

	// Test 3: Messages with whitespace should be trimmed
	os.Setenv(AuthMessagesEnvVar, " auth-token-4 , auth-token-5 ")
	token4 := generateTestToken(secret, "auth-token-4")
	token5 := generateTestToken(secret, "auth-token-5")

	if !ValidateToken(token4) {
		t.Error("ValidateToken() should return true for token4 with trimmed whitespace")
	}
	if !ValidateToken(token5) {
		t.Error("ValidateToken() should return true for token5 with trimmed whitespace")
	}

	// Test 4: Empty messages should be ignored
	os.Setenv(AuthMessagesEnvVar, "auth-token-6,,auth-token-7,  ,auth-token-8")
	token6 := generateTestToken(secret, "auth-token-6")
	token7 := generateTestToken(secret, "auth-token-7")
	token8 := generateTestToken(secret, "auth-token-8")

	if !ValidateToken(token6) {
		t.Error("ValidateToken() should return true for token6")
	}
	if !ValidateToken(token7) {
		t.Error("ValidateToken() should return true for token7")
	}
	if !ValidateToken(token8) {
		t.Error("ValidateToken() should return true for token8")
	}
}
