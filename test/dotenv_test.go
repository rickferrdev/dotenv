package dotenv_test

import (
	"os"
	"strings"
	"testing"

	"github.com/rickferrdev/dotenv"
)

type ConfigTest struct {
	Host      string  `env:"TEST_HOST" required:"true"`
	Port      int     `env:"TEST_PORT" required:"true"`
	Debug     bool    `env:"TEST_DEBUG" required:"true"`
	RateLimit float64 `env:"TEST_RATE" required:"true"`
	Ignored   string
}

type ConfigWithDefault struct {
	Auth string `env:"TEST_AUTH" required:"true" default:"xxx"`
	Name string `env:"TEST_NAME" default:"guest"`
}

type OptionalConfig struct {
	Name  string `env:"TEST_OPTIONAL_NAME"`
	Port  int    `env:"TEST_OPTIONAL_PORT"`
	Debug bool   `env:"TEST_OPTIONAL_DEBUG"`
}

func TestCollect(t *testing.T) {
	t.Setenv("TEST_VAR_1", "")
	t.Setenv("TEST_VAR_2", "")
	t.Setenv("TEST_VAR_3", "")
	t.Setenv("TEST_QUOTED", "")
	t.Setenv("TEST_SINGLE_QUOTED", "")
	t.Setenv("TEST_WITH_COMMENT", "")

	content := `
# Comment should be ignored
TEST_VAR_1=hello
export TEST_VAR_2=world
TEST_VAR_3=123
TEST_QUOTED="quoted value"
TEST_SINGLE_QUOTED='single quoted value'
TEST_WITH_COMMENT="secret" # inline comment
`

	tmpFile, err := os.CreateTemp("", ".env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	originalFilenames := dotenv.FilenameVariables
	dotenv.FilenameVariables = []string{tmpFile.Name()}
	defer func() {
		dotenv.FilenameVariables = originalFilenames
	}()

	dotenv.Collect()

	tests := map[string]string{
		"TEST_VAR_1":         "hello",
		"TEST_VAR_2":         "world",
		"TEST_VAR_3":         "123",
		"TEST_QUOTED":        "quoted value",
		"TEST_SINGLE_QUOTED": "single quoted value",
		"TEST_WITH_COMMENT":  "secret",
	}

	for key, expected := range tests {
		if got := os.Getenv(key); got != expected {
			t.Errorf("%s: expected %q, got %q", key, expected, got)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv("TEST_HOST", "localhost")
		t.Setenv("TEST_PORT", "8080")
		t.Setenv("TEST_DEBUG", "true")
		t.Setenv("TEST_RATE", "1.5")

		var cfg ConfigTest

		if err := dotenv.Unmarshal(&cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Host != "localhost" {
			t.Errorf("Host: expected %q, got %q", "localhost", cfg.Host)
		}

		if cfg.Port != 8080 {
			t.Errorf("Port: expected %d, got %d", 8080, cfg.Port)
		}

		if cfg.Debug != true {
			t.Errorf("Debug: expected true, got false")
		}

		if cfg.RateLimit != 1.5 {
			t.Errorf("RateLimit: expected %f, got %f", 1.5, cfg.RateLimit)
		}

		if cfg.Ignored != "" {
			t.Errorf("Ignored should not be filled, got %q", cfg.Ignored)
		}
	})

	t.Run("uses default when env is missing", func(t *testing.T) {
		os.Unsetenv("TEST_AUTH")
		os.Unsetenv("TEST_NAME")

		var cfg ConfigWithDefault

		if err := dotenv.Unmarshal(&cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Auth != "xxx" {
			t.Errorf("Auth: expected %q, got %q", "xxx", cfg.Auth)
		}

		if cfg.Name != "guest" {
			t.Errorf("Name: expected %q, got %q", "guest", cfg.Name)
		}
	})

	t.Run("uses default when env is empty", func(t *testing.T) {
		t.Setenv("TEST_AUTH", "")

		var cfg ConfigWithDefault

		if err := dotenv.Unmarshal(&cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Auth != "xxx" {
			t.Errorf("Auth: expected %q, got %q", "xxx", cfg.Auth)
		}
	})

	t.Run("required without env and default returns error", func(t *testing.T) {
		os.Unsetenv("TEST_HOST")
		os.Unsetenv("TEST_PORT")
		os.Unsetenv("TEST_DEBUG")
		os.Unsetenv("TEST_RATE")

		var cfg ConfigTest

		if err := dotenv.Unmarshal(&cfg); err == nil {
			t.Fatal("expected required error, got nil")
		}
	})

	t.Run("skips optional fields when env is missing", func(t *testing.T) {
		os.Unsetenv("TEST_OPTIONAL_NAME")
		os.Unsetenv("TEST_OPTIONAL_PORT")
		os.Unsetenv("TEST_OPTIONAL_DEBUG")

		cfg := OptionalConfig{
			Name:  "existing",
			Port:  3000,
			Debug: true,
		}

		if err := dotenv.Unmarshal(&cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Name != "existing" || cfg.Port != 3000 || cfg.Debug != true {
			t.Errorf("optional fields should be unchanged, got %+v", cfg)
		}
	})

	t.Run("non pointer returns error", func(t *testing.T) {
		var cfg ConfigTest

		if err := dotenv.Unmarshal(cfg); err == nil {
			t.Fatal("expected error for non-pointer value, got nil")
		}
	})

	t.Run("nil pointer returns error", func(t *testing.T) {
		var cfg *ConfigTest

		if err := dotenv.Unmarshal(cfg); err == nil {
			t.Fatal("expected error for nil pointer, got nil")
		}
	})

	t.Run("non struct pointer returns error", func(t *testing.T) {
		var value string

		if err := dotenv.Unmarshal(&value); err == nil {
			t.Fatal("expected error for pointer to non-struct, got nil")
		}
	})

	t.Run("invalid int returns error", func(t *testing.T) {
		t.Setenv("TEST_HOST", "localhost")
		t.Setenv("TEST_PORT", "not-a-number")
		t.Setenv("TEST_DEBUG", "true")
		t.Setenv("TEST_RATE", "1.5")

		var cfg ConfigTest

		if err := dotenv.Unmarshal(&cfg); err == nil {
			t.Fatal("expected int parse error, got nil")
		}
	})
}

func TestMarshal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := ConfigTest{
			Host:      "api.prod.com",
			Port:      9000,
			Debug:     false,
			RateLimit: 50.5,
			Ignored:   "This should not appear",
		}

		data, err := dotenv.Marshal(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := string(data)

		expected := []string{
			"TEST_HOST=api.prod.com",
			"TEST_PORT=9000",
			"TEST_DEBUG=false",
			"TEST_RATE=50.5",
		}

		for _, item := range expected {
			if !strings.Contains(output, item) {
				t.Errorf("expected output to contain %q, got:\n%s", item, output)
			}
		}

		if strings.Contains(output, "Ignored=") || strings.Contains(output, "This should not appear") {
			t.Errorf("field without env tag should be ignored, got:\n%s", output)
		}
	})

	t.Run("accepts struct value", func(t *testing.T) {
		cfg := ConfigTest{
			Host:      "api.prod.com",
			Port:      9000,
			Debug:     true,
			RateLimit: 50.5,
		}

		if _, err := dotenv.Marshal(cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("uses default when string field is empty", func(t *testing.T) {
		cfg := ConfigWithDefault{}

		data, err := dotenv.Marshal(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := string(data)

		if !strings.Contains(output, "TEST_AUTH=xxx") {
			t.Errorf("expected default auth, got:\n%s", output)
		}

		if !strings.Contains(output, "TEST_NAME=guest") {
			t.Errorf("expected default name, got:\n%s", output)
		}
	})

	t.Run("quotes values with spaces", func(t *testing.T) {
		cfg := struct {
			Name string `env:"TEST_NAME"`
		}{
			Name: "john doe",
		}

		data, err := dotenv.Marshal(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := `TEST_NAME="john doe"`
		if !strings.Contains(string(data), expected) {
			t.Errorf("expected %q, got:\n%s", expected, string(data))
		}
	})

	t.Run("required empty string without default returns error", func(t *testing.T) {
		cfg := struct {
			Token string `env:"TEST_TOKEN" required:"true"`
		}{}

		if _, err := dotenv.Marshal(&cfg); err == nil {
			t.Fatal("expected required error, got nil")
		}
	})

	t.Run("nil pointer returns error", func(t *testing.T) {
		var cfg *ConfigTest

		if _, err := dotenv.Marshal(cfg); err == nil {
			t.Fatal("expected error for nil pointer, got nil")
		}
	})

	t.Run("non struct returns error", func(t *testing.T) {
		if _, err := dotenv.Marshal("invalid"); err == nil {
			t.Fatal("expected error for non-struct, got nil")
		}
	})
}
