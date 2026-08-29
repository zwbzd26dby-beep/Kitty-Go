// Package deepseek provides the DeepSeek provider.
package deepseek

import (
	"net/http"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
)

// DefaultBaseURL is the DeepSeek API base URL.
const DefaultBaseURL = "https://api.deepseek.com/v1"

// EnvKey is the environment variable holding the DeepSeek API key.
const EnvKey = "DEEPSEEK_API_KEY"

// Options configures the DeepSeek provider.
type Options struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
	Limits  provider.Limits
}

// New creates a DeepSeek provider.
func New(opts Options) provider.Provider {
	base := opts.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	return provider.NewOpenAICompat(provider.Options{
		Name:    "deepseek",
		APIKey:  opts.APIKey,
		BaseURL: base,
		Client:  opts.Client,
		Limits:  opts.Limits,
	})
}
