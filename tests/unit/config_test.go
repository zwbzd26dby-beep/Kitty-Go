package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/config"
)

func TestConfigDefaults(t *testing.T) {
	cfg := config.Default()
	if cfg.Provider != "mock" {
		t.Fatalf("expected provider mock, got %q", cfg.Provider)
	}
	if cfg.Model != "mock-model" {
		t.Fatalf("expected model mock-model, got %q", cfg.Model)
	}
	if cfg.Timeout <= 0 {
		t.Fatalf("expected positive timeout, got %v", cfg.Timeout)
	}
	if cfg.REPL.SecurityLevel != "standard" {
		t.Fatalf("expected security_level standard, got %q", cfg.REPL.SecurityLevel)
	}
}

func TestParseUserOverridesDefaults(t *testing.T) {
	data := []byte(`
provider: opencode
model: big-pickle
api_key_env: OPENCODE_API_KEY
opencode:
  base_url: "https://opencode.ai/zen/v1"
  default_model: "big-pickle"
budget:
  daily: 5.0
  monthly: 50.0
limits:
  requests_daily: 1000
repl:
  security_level: high
`)
	cfg, err := config.ParseUser(data)
	if err != nil {
		t.Fatalf("ParseUser: %v", err)
	}
	if cfg.Provider != "opencode" {
		t.Fatalf("expected opencode, got %q", cfg.Provider)
	}
	if cfg.Model != "big-pickle" {
		t.Fatalf("expected big-pickle, got %q", cfg.Model)
	}
	if cfg.Budget.Daily != 5.0 {
		t.Fatalf("expected daily budget 5.0, got %v", cfg.Budget.Daily)
	}
	if cfg.REPL.SecurityLevel != "high" {
		t.Fatalf("expected high, got %q", cfg.REPL.SecurityLevel)
	}
	oc, ok := cfg.Providers["opencode"]
	if !ok {
		t.Fatal("expected opencode provider settings")
	}
	if oc.BaseURL != "https://opencode.ai/zen/v1" {
		t.Fatalf("unexpected opencode base url %q", oc.BaseURL)
	}
}

func TestLoadFromTempUserFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("provider: opencode\nmodel: big-pickle\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(config.LoadOptions{UserPath: path})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Provider != "opencode" {
		t.Fatalf("expected opencode, got %q", cfg.Provider)
	}
	if cfg.Model != "big-pickle" {
		t.Fatalf("expected big-pickle, got %q", cfg.Model)
	}
	// Defaults survive when user config is partial.
	if cfg.REPL.SecurityLevel != "standard" {
		t.Fatalf("expected default security_level, got %q", cfg.REPL.SecurityLevel)
	}
}

func TestLoadMissingUserFileNotError(t *testing.T) {
	cfg, err := config.Load(config.LoadOptions{UserPath: filepath.Join(t.TempDir(), "nope.yaml")})
	if err != nil {
		t.Fatalf("expected missing user file to be non-fatal: %v", err)
	}
	if cfg.Provider != "mock" {
		t.Fatalf("expected default provider, got %q", cfg.Provider)
	}
}

func TestLoadValidationMessages(t *testing.T) {
	data := []byte("repl:\n  security_level: bogus\n")
	cfg, err := config.ParseUser(data)
	if err != nil {
		t.Fatal(err)
	}
	msgs := config.Validate(cfg)
	found := false
	for _, m := range msgs {
		if m.Level == "warning" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected a validation warning")
	}
}
