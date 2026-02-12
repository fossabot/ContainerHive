package secrets

import (
	"testing"
)

func TestPlainTextResolver_Resolve(t *testing.T) {
	resolver := &PlainTextResolver{}

	tests := []struct {
		name          string
		value         string
		expectedValue string
		expectError   bool
	}{
		{
			name:          "empty string",
			value:         "",
			expectedValue: "",
			expectError:   false,
		},
		{
			name:          "simple text",
			value:         "my-secret-value",
			expectedValue: "my-secret-value",
			expectError:   false,
		},
		{
			name:          "text with spaces",
			value:         "my secret value",
			expectedValue: "my secret value",
			expectError:   false,
		},
		{
			name:          "text with special characters",
			value:         "my-secret!@#$%^&*()",
			expectedValue: "my-secret!@#$%^&*()",
			expectError:   false,
		},
		{
			name:          "multi-line text",
			value:         "line1\nline2\nline3",
			expectedValue: "line1\nline2\nline3",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := resolver.Resolve(tt.value)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if value != tt.expectedValue {
				t.Errorf("expected value %q, got %q", tt.expectedValue, value)
			}
		})
	}
}

func TestEnvVarResolver_Resolve(t *testing.T) {
	resolver := &EnvVarResolver{}

	// Set up test environment variables
	t.Setenv("TEST_SECRET", "test-value")
	t.Setenv("ANOTHER_VAR", "another-value")
	t.Setenv("EMPTY_VAR", "")

	tests := []struct {
		name          string
		value         string
		expectedValue string
		expectError   bool
		errorContains string
	}{
		{
			name:          "valid env var with braces",
			value:         "${TEST_SECRET}",
			expectedValue: "test-value",
			expectError:   false,
		},
		{
			name:          "valid env var without braces",
			value:         "$TEST_SECRET",
			expectedValue: "test-value",
			expectError:   false,
		},
		{
			name:          "another valid env var",
			value:         "${ANOTHER_VAR}",
			expectedValue: "another-value",
			expectError:   false,
		},
		{
			name:          "empty env var",
			value:         "${EMPTY_VAR}",
			expectedValue: "",
			expectError:   false,
		},
		{
			name:          "non-existent env var",
			value:         "${NON_EXISTENT}",
			expectedValue: "",
			expectError:   true,
			errorContains: "not found",
		},
		{
			name:          "invalid env var format",
			value:         "not-an-env-var",
			expectedValue: "",
			expectError:   false, // Should return empty string for non-matching format
		},
		{
			name:          "env var with invalid name",
			value:         "${123INVALID}",
			expectedValue: "",
			expectError:   false, // Should return empty string for invalid names
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := resolver.Resolve(tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorContains != "" && !containsError(err, tt.errorContains) {
					t.Errorf("expected error containing %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if value != tt.expectedValue {
				t.Errorf("expected value %q, got %q", tt.expectedValue, value)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	// Set up test environment variables
	t.Setenv("TEST_ENV_SECRET", "env-secret-value")

	tests := []struct {
		name          string
		secretType    string
		value         string
		expectedValue string
		expectError   bool
		errorContains string
	}{
		{
			name:          "plain text secret with explicit type",
			secretType:    "plain",
			value:         "plain-secret-value",
			expectedValue: "plain-secret-value",
			expectError:   false,
		},
		{
			name:          "env var secret with explicit type",
			secretType:    "env",
			value:         "${TEST_ENV_SECRET}",
			expectedValue: "env-secret-value",
			expectError:   false,
		},
		{
			name:          "plain text secret auto-detected",
			secretType:    "",
			value:         "auto-detected-plain",
			expectedValue: "auto-detected-plain",
			expectError:   false,
		},
		{
			name:          "env var secret auto-detected",
			secretType:    "",
			value:         "${TEST_ENV_SECRET}",
			expectedValue: "env-secret-value",
			expectError:   false,
		},
		{
			name:          "non-existent env var auto-detected",
			secretType:    "",
			value:         "${NON_EXISTENT_VAR}",
			expectedValue: "",
			expectError:   true,
			errorContains: "not found",
		},
		{
			name:          "invalid secret type",
			secretType:    "invalid",
			value:         "some-value",
			expectedValue: "",
			expectError:   true,
			errorContains: "no resolver could handle",
		},
		{
			name:          "empty value with plain type",
			secretType:    "plain",
			value:         "",
			expectedValue: "",
			expectError:   false,
		},
		{
			name:          "empty value auto-detected",
			secretType:    "",
			value:         "",
			expectedValue: "",
			expectError:   true, // Empty string can't be auto-detected
			errorContains: "no resolver could handle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := Resolve(tt.secretType, tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorContains != "" && !containsError(err, tt.errorContains) {
					t.Errorf("expected error containing %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if value != tt.expectedValue {
				t.Errorf("expected value %q, got %q", tt.expectedValue, value)
			}
		})
	}
}

// Helper function to check if error message contains expected substring
func containsError(err error, substring string) bool {
	return err != nil && containsString(err.Error(), substring)
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || containsString(s[1:], substr)))
}
