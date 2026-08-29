package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// userYAML mirrors the user-level config file (configs/user.yaml or
// ~/.kitty/config.yaml). Optional per-provider blocks configure base URLs and
// default models without requiring secrets.
type userYAML struct {
	Provider  string  `yaml:"provider"`
	Model     string  `yaml:"model"`
	APIKey    string  `yaml:"api_key"`
	APIKeyEnv string  `yaml:"api_key_env"`
	BaseURL   string  `yaml:"base_url"`
	Timeout   float64 `yaml:"timeout"`

	OpenCode   *provBlock `yaml:"opencode"`
	OpenAI     *provBlock `yaml:"openai"`
	Kimi       *provBlock `yaml:"kimi"`
	DeepSeek   *provBlock `yaml:"deepseek"`
	OpenRouter *provBlock `yaml:"openrouter"`
	Ollama     *provBlock `yaml:"ollama"`

	Limits limitsYAML `yaml:"limits"`
	Budget budgetYAML `yaml:"budget"`
	REPL   replYAML   `yaml:"repl"`
}

type provBlock struct {
	BaseURL      string `yaml:"base_url"`
	DefaultModel string `yaml:"default_model"`
	MaxTokens    int    `yaml:"max_tokens"`
}

type limitsYAML struct {
	MaxRequestsPerMinute uint32 `yaml:"max_requests_per_minute"`
	MaxTokensPerMinute   uint32 `yaml:"max_tokens_per_minute"`
	TokensMonthly        uint64 `yaml:"tokens_monthly"`
	RequestsDaily        uint32 `yaml:"requests_daily"`
}

type budgetYAML struct {
	Daily   float64 `yaml:"daily"`
	Monthly float64 `yaml:"monthly"`
}

type replYAML struct {
	SecurityLevel string `yaml:"security_level"`
	ShowCost      *bool  `yaml:"show_cost"`
	SandboxRoot   string `yaml:"sandbox_root"`
}

// ParseUser parses a user-level YAML config into a Config, overlaying it on
// top of the defaults.
func ParseUser(data []byte) (*Config, error) {
	u := userYAML{}
	if err := yaml.Unmarshal(data, &u); err != nil {
		return nil, fmt.Errorf("parse user config: %w", err)
	}
	cfg := Default()
	cfg.applyUser(u)
	return cfg, nil
}

func (c *Config) applyUser(u userYAML) {
	if u.Provider != "" {
		c.Provider = u.Provider
	}
	if u.Model != "" {
		c.Model = u.Model
	}
	if u.APIKey != "" {
		c.APIKey = u.APIKey
	}
	if u.APIKeyEnv != "" {
		c.APIKeyEnv = u.APIKeyEnv
	}
	if u.BaseURL != "" {
		c.BaseURL = u.BaseURL
	}
	if u.Timeout > 0 {
		c.Timeout = u.Timeout
	}

	setBlock := func(name string, b *provBlock) {
		if b == nil {
			return
		}
		p := c.Providers[name]
		if b.BaseURL != "" {
			p.BaseURL = b.BaseURL
		}
		if b.DefaultModel != "" {
			p.Model = b.DefaultModel
		}
		if b.MaxTokens > 0 {
			p.MaxTokens = b.MaxTokens
		}
		c.Providers[name] = p
	}
	setBlock("opencode", u.OpenCode)
	setBlock("openai", u.OpenAI)
	setBlock("kimi", u.Kimi)
	setBlock("deepseek", u.DeepSeek)
	setBlock("openrouter", u.OpenRouter)
	setBlock("ollama", u.Ollama)

	c.Limits = LimitsConfig{
		MaxRequestsPerMinute: u.Limits.MaxRequestsPerMinute,
		MaxTokensPerMinute:   u.Limits.MaxTokensPerMinute,
		TokensMonthly:        u.Limits.TokensMonthly,
		RequestsDaily:        u.Limits.RequestsDaily,
	}
	c.Budget = BudgetConfig{Daily: u.Budget.Daily, Monthly: u.Budget.Monthly}
	if u.REPL.SecurityLevel != "" {
		c.REPL.SecurityLevel = u.REPL.SecurityLevel
	}
	if u.REPL.ShowCost != nil {
		c.REPL.ShowCost = *u.REPL.ShowCost
	}
	if u.REPL.SandboxRoot != "" {
		c.REPL.SandboxRoot = u.REPL.SandboxRoot
	}
}
