package secrets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVaultResolver_Resolve(t *testing.T) {
	resolver := &VaultSecretResolver{}

	tests := []struct {
		name          string
		value         string
		expectedValue string
		expectError   bool
		errorContains string
		setupEnv      func()
		cleanupEnv    func()
	}{
		{
			name:          "non-vault value",
			value:         "plain-text-secret",
			expectedValue: "",
			expectError:   false,
		},
		{
			name:          "vault value without hash",
			value:         "vault://secret/data/myapp",
			expectedValue: "",
			expectError:   true,
			errorContains: "malformed vault secret spec",
		},
		{
			name:          "vault value with empty path",
			value:         "vault://#field",
			expectedValue: "",
			expectError:   true,
			errorContains: "path cannot be empty",
		},
		{
			name:          "vault value with empty field",
			value:         "vault://secret/data/myapp#",
			expectedValue: "",
			expectError:   true,
			errorContains: "field cannot be empty",
		},
		{
			name:          "vault value with missing token",
			value:         "vault://secret/data/myapp#password",
			expectedValue: "",
			expectError:   true,
			errorContains: "open /tmp/.vault-token: no such file or directory",
			setupEnv: func() {
				os.Unsetenv("VAULT_ADDR")
				os.Unsetenv("VAULT_TOKEN")
				os.Setenv("HOME", "/tmp")
			},
			cleanupEnv: func() {
				os.Unsetenv("HOME")
			},
		},
		{
			name:          "vault value with special characters",
			value:         "vault://secret/data/app-with-dashes#api_key",
			expectedValue: "",
			expectError:   true,
			errorContains: "$HOME is not defined",
		},
		{
			name:          "vault value with complex path",
			value:         "vault://secret/data/nested/deep/path#token",
			expectedValue: "",
			expectError:   true,
			errorContains: "$HOME is not defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment if needed
			if tt.setupEnv != nil {
				defer tt.cleanupEnv()
				tt.setupEnv()
			}

			value, err := resolver.Resolve(tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("Resolve() error = nil, want non-nil")
				} else if tt.errorContains != "" && (err == nil || !containsString(err.Error(), tt.errorContains)) {
					t.Errorf("Resolve() error = %v, want error containing %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("Resolve() error = %v, want nil", err)
				}
			}

			if value != tt.expectedValue {
				t.Errorf("Resolve() = %v, want %v", value, tt.expectedValue)
			}
		})
	}
}

func TestVaultResolver_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	resolver := &VaultSecretResolver{}

	// Set up test environment
	tempDir := t.TempDir()
	vaultTokenFile := filepath.Join(tempDir, ".vault-token")
	err := os.WriteFile(vaultTokenFile, []byte("test-root-token"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test vault token file: %v", err)
	}

	t.Setenv("HOME", tempDir)
	t.Setenv("VAULT_ADDR", "http://localhost:8200")
	t.Setenv("VAULT_TOKEN", "test-root-token")

	tests := []struct {
		name          string
		value         string
		expectError   bool
		errorContains string
	}{
		{
			name:          "vault secret with non-existent path",
			value:         "vault://secret/data/nonexistent#password",
			expectError:   true,
			errorContains: "dial tcp",
		},
		{
			name:          "vault secret with malformed path",
			value:         "vault://invalid/path#field",
			expectError:   true,
			errorContains: "dial tcp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := resolver.Resolve(tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("Resolve() error = nil, want non-nil")
				} else if tt.errorContains != "" && (err == nil || !containsString(err.Error(), tt.errorContains)) {
					t.Errorf("Resolve() error = %v, want error containing %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("Resolve() error = %v, want nil", err)
				}
			}

			// Should return empty string for errors
			if !tt.expectError && value != "" {
				t.Errorf("Resolve() = %v, want empty string", value)
			}
		})
	}

}

func TestVaultResolver_DirectResolution(t *testing.T) {
	resolver := &VaultSecretResolver{}

	tests := []struct {
		name          string
		value         string
		expectError   bool
		errorContains string
	}{
		{
			name:          "vault value with malformed spec",
			value:         "vault://missing-hash",
			expectError:   true,
			errorContains: "malformed vault secret spec",
		},
		{
			name:          "vault value with correct format but missing config",
			value:         "vault://secret/data/app#password",
			expectError:   true,
			errorContains: "$HOME is not defined",
		},
		{
			name:        "non-vault value should return empty",
			value:       "plain-text-secret",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := resolver.Resolve(tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("Resolve() error = nil, want non-nil")
				} else if tt.errorContains != "" && (err == nil || !strings.Contains(err.Error(), tt.errorContains)) {
					t.Errorf("Resolve() error = %v, want error containing %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("Resolve() error = %v, want nil", err)
				}
			}

			// Should return empty string for non-vault values or errors
			if value != "" {
				t.Errorf("Resolve() = %v, want empty string", value)
			}
		})
	}
}
