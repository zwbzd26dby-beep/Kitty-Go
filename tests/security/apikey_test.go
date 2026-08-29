package securitytest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/config"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/security"
)

func TestGetAPIKeyFromSeed(t *testing.T) {
	m := security.NewManager(map[string]string{"opencode": "seed-key"}, security.DefaultEnvFallback())
	got, err := m.GetAPIKey("opencode")
	if err != nil {
		t.Fatalf("GetAPIKey: %v", err)
	}
	if got != "seed-key" {
		t.Fatalf("expected seed-key, got %q", got)
	}
}

func TestGetAPIKeyFromEnvFallback(t *testing.T) {
	t.Setenv("OPENCODE_API_KEY", "env-key")
	m := security.NewManager(nil, security.DefaultEnvFallback())
	got, err := m.GetAPIKey("opencode")
	if err != nil {
		t.Fatalf("GetAPIKey: %v", err)
	}
	if got != "env-key" {
		t.Fatalf("expected env-key, got %q", got)
	}
}

func TestSetAPIKeyOverrides(t *testing.T) {
	t.Setenv("OPENCODE_API_KEY", "env-key")
	m := security.NewManager(nil, security.DefaultEnvFallback())
	m.SetAPIKey("opencode", "memory-key")
	got, err := m.GetAPIKey("opencode")
	if err != nil {
		t.Fatal(err)
	}
	if got != "memory-key" {
		t.Fatalf("expected memory-key, got %q", got)
	}
}

func TestMissingAPIKeyErrors(t *testing.T) {
	m := security.NewManager(nil, map[string]string{})
	_, err := m.GetAPIKey("unknown-provider")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if m.HasAPIKey("unknown-provider") {
		t.Fatal("HasAPIKey should be false")
	}
}

func TestLoadSecretsThroughConfig(t *testing.T) {
	dir := t.TempDir()
	secrets := filepath.Join(dir, "secrets.yaml")
	if err := os.WriteFile(secrets, []byte("secrets:\n  provider_api_keys:\n    opencode: \"file-key\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(config.LoadOptions{SecretPath: secrets})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	got, err := cfg.Security.GetAPIKey("opencode")
	if err != nil {
		t.Fatalf("GetAPIKey: %v", err)
	}
	if got != "file-key" {
		t.Fatalf("expected file-key, got %q", got)
	}
}
