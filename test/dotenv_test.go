package dotenv_test

import (
	"os"
	"strings"
	"testing"

	"github.com/rickferrdev/dotenv"
)

// ConfigTest is a sample structure for Marshal/Unmarshal tests
type ConfigTest struct {
	Host      string  `env:"TEST_HOST"`
	Port      int     `env:"TEST_PORT"`
	Debug     bool    `env:"TEST_DEBUG"`
	RateLimit float64 `env:"TEST_RATE"`
	Ignored   string  // Field without tag, should be ignored
}

// TestCollect verifies if the file is read and variables are injected into the environment
func TestCollect(t *testing.T) {
	// 1. Setup: Create a temporary .env file
	// We include quoted values here to test the 'utils' logic implicitly
	content := `
# Comment should be ignored
TEST_VAR_1=hello
export TEST_VAR_2=world

# Empty line above
TEST_VAR_3=123

# Testing the internal quotes utility
TEST_QUOTED="quoted value"
TEST_SINGLE_QUOTED='single quoted value'
TEST_WITH_COMMENT="secret" # inline comment
`
	tmpFile, err := os.CreateTemp("", ".env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file at the end

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// 2. Mock: Temporarily replace the target file list
	// to point to our temporary file
	originalFilenames := dotenv.FilenameVariables
	dotenv.FilenameVariables = []string{tmpFile.Name()}

	// Ensure we revert to the original filenames after the test
	defer func() { dotenv.FilenameVariables = originalFilenames }()

	// 3. Execution
	dotenv.Collect()

	// 4. Assertions
	tests := []struct {
		key      string
		expected string
	}{
		{"TEST_VAR_1", "hello"},
		{"TEST_VAR_2", "world"},
		{"TEST_VAR_3", "123"},
		{"TEST_QUOTED", "quoted value"},               // Validates utils.quotes
		{"TEST_SINGLE_QUOTED", "single quoted value"}, // Validates utils.quotes
		{"TEST_WITH_COMMENT", "secret"},               // Validates utils.quotes
	}

	for _, tt := range tests {
		val := os.Getenv(tt.key)
		if val != tt.expected {
			t.Errorf("Collect(): for key %s, expected '%s', got '%s'", tt.key, tt.expected, val)
		}
		// Clean up created env vars to avoid polluting other tests
		os.Unsetenv(tt.key)
	}
}

// TestUnmarshal verifies if environment variables correctly fill the struct
func TestUnmarshal(t *testing.T) {
	// 1. Setup: Define real environment variables
	envVars := map[string]string{
		"TEST_HOST":  "localhost",
		"TEST_PORT":  "8080",
		"TEST_DEBUG": "true",
		"TEST_RATE":  "1.5",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
		defer os.Unsetenv(k) // Clean up at the end
	}

	// 2. Success Test
	t.Run("Success", func(t *testing.T) {
		var cfg ConfigTest
		err := dotenv.Unmarshal(&cfg)

		if err != nil {
			t.Fatalf("Unmarshal returned unexpected error: %v", err)
		}

		if cfg.Host != "localhost" {
			t.Errorf("Host: expected 'localhost', got '%s'", cfg.Host)
		}
		if cfg.Port != 8080 {
			t.Errorf("Port: expected 8080, got %d", cfg.Port)
		}
		if !cfg.Debug {
			t.Errorf("Debug: expected true, got false")
		}
		if cfg.RateLimit != 1.5 {
			t.Errorf("RateLimit: expected 1.5, got %f", cfg.RateLimit)
		}
	})

	// 3. Error Tests
	t.Run("Errors", func(t *testing.T) {
		// Case 1: Pass non-pointer value
		var cfg ConfigTest
		err := dotenv.Unmarshal(cfg)
		if err == nil {
			t.Error("Expected error when passing struct by value, but got nil")
		}

		// Case 2: Type error (trying to parse text into int)
		os.Setenv("TEST_PORT", "not-a-number")
		err = dotenv.Unmarshal(&cfg)
		if err == nil {
			t.Error("Expected int parse error, got nil")
		}
	})
}

// TestMarshal verifies if the struct is correctly converted to .env format
func TestMarshal(t *testing.T) {
	cfg := ConfigTest{
		Host:      "api.prod.com",
		Port:      9000,
		Debug:     false,
		RateLimit: 50.5,
		Ignored:   "This should not appear",
	}

	bytes, err := dotenv.Marshal(&cfg)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	output := string(bytes)

	// Verifications
	expectedContains := []string{
		"TEST_HOST=api.prod.com",
		"TEST_PORT=9000",
		"TEST_DEBUG=false",
		"TEST_RATE=50.5",
	}

	for _, exp := range expectedContains {
		if !strings.Contains(output, exp) {
			t.Errorf("Marshal output does not contain: %s. Full output:\n%s", exp, output)
		}
	}

	// Verify if the field without tag was ignored (should not have empty key or field name)
	if strings.Contains(output, "Ignored=") || strings.Contains(output, "=This should not appear") {
		t.Error("Marshal included a field that should not have an env tag")
	}
}
