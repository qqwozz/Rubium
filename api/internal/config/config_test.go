package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "resolved")
	defer os.Unsetenv("TEST_VAR")

	tests := []struct {
		input    string
		expected string
	}{
		{"${TEST_VAR}", "resolved"},
		{"plain", "plain"},
		{"${MISSING}", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := resolveEnv(tt.input)
			if got != tt.expected {
				t.Errorf("resolveEnv(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLoad_UsesEnvPort(t *testing.T) {
	dir := t.TempDir()
	yamlFile := filepath.Join(dir, "test_config.yaml")
	content := `
supabase:
  url: http://localhost
  anon_key: test-anon
  service_key: test-service
server:
  go_port: 9000
`
	if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	os.Setenv("CONFIG_YAML", yamlFile)
	os.Setenv("PORT", "7777")
	defer func() {
		os.Unsetenv("CONFIG_YAML")
		os.Unsetenv("PORT")
	}()

	cfg := Load()
	if cfg.Port != "7777" {
		t.Errorf("expected PORT env to override, got %s", cfg.Port)
	}
	if cfg.SupabaseURL != "http://localhost" {
		t.Errorf("unexpected url: %s", cfg.SupabaseURL)
	}
}