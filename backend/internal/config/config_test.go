// config/config_test.go
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	originalApiKey := os.Getenv("GEMINI_API_KEY")
	originalPort := os.Getenv("PORT")
	defer func() {
		os.Setenv("GEMINI_API_KEY", originalApiKey)
		os.Setenv("PORT", originalPort)
	}()

	testCases := []struct {
		name         string
		setEnv       func(t *testing.T)
		expectError  bool
		expectedPort string
		expectedKey  string
	}{
		{
			name: "Valid config",
			setEnv: func(t *testing.T) {
				t.Setenv("GEMINI_API_KEY", "test-api-key")
				t.Setenv("PORT", "9090")
			},
			expectError:  false,
			expectedPort: "9090",
			expectedKey:  "test-api-key",
		},
		{
			name: "API Key missing",
			setEnv: func(t *testing.T) {
				t.Setenv("GEMINI_API_KEY", "")
				t.Setenv("PORT", "9090")
			},
			expectError:  true,
			expectedPort: "",
			expectedKey:  "",
		},
		{
			name: "Port missing (uses default)",
			setEnv: func(t *testing.T) {
				t.Setenv("GEMINI_API_KEY", "test-api-key-2")
				os.Unsetenv("PORT")
			},
			expectError:  false,
			expectedPort: "8080",
			expectedKey:  "test-api-key-2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setEnv(t)

			cfg, err := LoadConfig()

			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none")
				assert.Nil(t, cfg, "Config should be nil on error")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				require.NotNil(t, cfg, "Config should not be nil on success")
				assert.Equal(t, tc.expectedKey, cfg.GeminiAPIKey)
				assert.Equal(t, tc.expectedPort, cfg.ServerPort)
				assert.Equal(t, []string{"http://localhost:5173"}, cfg.AllowOrigins)
			}
		})
	}
}
