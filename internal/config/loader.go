package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/security"
)

// LoadOptions configures the config loader.
type LoadOptions struct {
	// UserPath overrides the user config file location. When empty, it falls
	// back to ~/.kitty/config.yaml if present.
	UserPath string
	// SecretPath overrides the secrets file location (fallback: environment).
	SecretPath string
}

// Load merges defaults, user config and secrets into a validated Config.
// Missing user/secret files are not errors.
func Load(opts LoadOptions) (*Config, error) {
	cfg := Default()

	if opts.UserPath == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			opts.UserPath = filepath.Join(home, ".kitty", "config.yaml")
		}
	}
	userData, err := readOptional(opts.UserPath)
	if err != nil {
		return nil, err
	}
	if userData != nil {
		parsed, err := ParseUser(userData)
		if err != nil {
			return nil, err
		}
		merged := cfg
		*cfg = *parsed
		_ = merged
	}

	messages := Validate(cfg)

	// Resolve explicitly-provided api_key (from user config) into the manager.
	seed := map[string]string{}
	if cfg.APIKey != "" {
		seed[cfg.Provider] = cfg.APIKey
	}
	if opts.SecretPath != "" {
		if s, serr := loadSecretsYAML(opts.SecretPath); serr == nil {
			for k, v := range s {
				seed[k] = v
			}
		}
	}
	cfg.Security = security.NewManager(seed, security.DefaultEnvFallback())

	cfg.messages = messages
	return cfg, nil
}

func readOptional(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// loadSecretsYAML parses a secrets.yaml file into a provider->key map.
func loadSecretsYAML(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc struct {
		Secrets struct {
			ProviderAPIKeys map[string]string `yaml:"provider_api_keys"`
		} `yaml:"secrets"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse secrets: %w", err)
	}
	if doc.Secrets.ProviderAPIKeys == nil {
		return map[string]string{}, nil
	}
	return doc.Secrets.ProviderAPIKeys, nil
}
