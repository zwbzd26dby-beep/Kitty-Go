// Package security centralises API-key and secrets handling
// (Master Architecture §22, expanded in Phase 14).
package security

import (
	"fmt"
	"os"
)

// Manager resolves and stores provider API keys. Keys are kept in memory and
// never logged or written to config files by this package.
type Manager struct {
	keys        map[string]string
	envFallback map[string]string
}

// NewManager creates a Manager seeded from an explicit provider->key map
// (typically loaded from secrets.yaml) and an env-fallback map.
func NewManager(seed map[string]string, envFallback map[string]string) *Manager {
	keys := make(map[string]string, len(seed))
	for k, v := range seed {
		keys[k] = v
	}
	return &Manager{keys: keys, envFallback: envFallback}
}

// DefaultEnvFallback returns the default provider->environment variable map.
func DefaultEnvFallback() map[string]string {
	return map[string]string{
		"openai":     "OPENAI_API_KEY",
		"kimi":       "KIMI_API_KEY",
		"deepseek":   "DEEPSEEK_API_KEY",
		"openrouter": "OPENROUTER_API_KEY",
		"ollama":     "OLLAMA_API_KEY",
		"opencode":   "OPENCODE_API_KEY",
	}
}

// GetAPIKey returns the effective key for a provider. Resolution order is:
// an explicitly set/loaded key first, then the configured environment
// variable. Returns "" if neither is present.
func (m *Manager) GetAPIKey(provider string) (string, error) {
	if k, ok := m.keys[provider]; ok && k != "" {
		return k, nil
	}
	if env, ok := m.envFallback[provider]; ok && env != "" {
		if v := os.Getenv(env); v != "" {
			return v, nil
		}
	}
	return "", fmt.Errorf("no API key for provider %q", provider)
}

// SetAPIKey stores a key in memory for the given provider.
func (m *Manager) SetAPIKey(provider, key string) {
	if m.keys == nil {
		m.keys = make(map[string]string)
	}
	m.keys[provider] = key
}

// HasAPIKey reports whether a key is resolvable for provider.
func (m *Manager) HasAPIKey(provider string) bool {
	_, err := m.GetAPIKey(provider)
	return err == nil
}
