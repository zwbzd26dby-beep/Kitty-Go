package config

// Default returns a Config populated with the plugin/startup defaults.
// These mirror configs/default.yaml and must never contain secrets.
func Default() *Config {
	return &Config{
		Provider:  "mock",
		Model:     "mock-model",
		APIKeyEnv: "OPENCODE_API_KEY",
		Timeout:   30.0,
		Providers: map[string]ProviderSettings{},
		Budget:    BudgetConfig{},
		REPL: REPLConfig{
			SecurityLevel: "standard",
			ShowCost:      true,
		},
	}
}
