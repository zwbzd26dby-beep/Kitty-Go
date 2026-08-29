// Package config implements the three-tier configuration system
// (default / user / secret) and the Config interface (Master Architecture §22).
package config

import "github.com/zwbzd26dby-beep/Kitty-Go/internal/security"

// Message is a side-channel configuration note (deprecations, hints) surfaced
// to the caller after loading.
type Message struct {
	Level   string // "info", "warning"
	Message string
}

// ProviderSettings configures a single provider.
type ProviderSettings struct {
	Provider  string        `yaml:"provider"`
	Model     string        `yaml:"model"`
	BaseURL   string        `yaml:"base_url"`
	APIKeyEnv string        `yaml:"api_key_env"`
	TimeoutMs int           `yaml:"timeout_ms"`
	Pricing   PricingConfig `yaml:"pricing"`
	MaxTokens int           `yaml:"max_tokens"`
}

// PricingConfig carries static model pricing for a provider.
type PricingConfig struct {
	Input  float64 `yaml:"input"`
	Output float64 `yaml:"output"`
}

// LimitsConfig configures rate and token limits.
type LimitsConfig struct {
	MaxRequestsPerMinute uint32 `yaml:"max_requests_per_minute"`
	MaxTokensPerMinute   uint32 `yaml:"max_tokens_per_minute"`
	TokensMonthly        uint64 `yaml:"tokens_monthly"`
	RequestsDaily        uint32 `yaml:"requests_daily"`
}

// BudgetConfig configures daily/monthly spend caps.
type BudgetConfig struct {
	Daily   float64 `yaml:"daily"`
	Monthly float64 `yaml:"monthly"`
}

// REPLConfig configures the REPL behaviour.
type REPLConfig struct {
	SecurityLevel string `yaml:"security_level"`
	ShowCost      bool   `yaml:"show_cost"`
	SandboxRoot   string `yaml:"sandbox_root"`
}

// Config is the merged, validated runtime configuration.
type Config struct {
	Provider      string  `yaml:"provider"`
	Model         string  `yaml:"model"`
	APIKey        string  `yaml:"api_key"`
	BaseURL       string  `yaml:"base_url"`
	APIKeyEnv     string  `yaml:"api_key_env"`
	Timeout       float64 `yaml:"timeout"`
	SecurityLevel string  `yaml:"security_level"`

	Providers map[string]ProviderSettings `yaml:"-"`
	Limits    LimitsConfig                `yaml:"limits"`
	Budget    BudgetConfig                `yaml:"budget"`
	REPL      REPLConfig                  `yaml:"repl"`

	// Security resolves provider API keys (not serialised).
	Security *security.Manager `yaml:"-"`
	// messages holds validation notes surfaced during load.
	messages []Message
}

// Messages returns non-fatal validation notes collected during Load.
func (c *Config) Messages() []Message { return c.messages }

// Settings exposes the public read-only configuration view (Master Arch §22).
type Settings interface {
	GetProvider() string
	GetModel() string
	GetLimits() LimitsConfig
	GetBudget() BudgetConfig
	GetREPL() REPLConfig
	Messages() []Message
}

var _ Settings = (*Config)(nil)

// GetProvider implements Settings.
func (c *Config) GetProvider() string { return c.Provider }

// GetModel implements Settings.
func (c *Config) GetModel() string { return c.Model }

// GetLimits implements Settings.
func (c *Config) GetLimits() LimitsConfig { return c.Limits }

// GetBudget implements Settings.
func (c *Config) GetBudget() BudgetConfig { return c.Budget }

// GetREPL implements Settings.
func (c *Config) GetREPL() REPLConfig { return c.REPL }
