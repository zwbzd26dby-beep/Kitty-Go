// Package kimi provides the Kimi (Moonshot) provider.
package kimi

import (
	"net/http"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
)

// DefaultBaseURL is the Kimi (Moonshot) API base URL.
const DefaultBaseURL = "https://api.moonshot.cn/v1"

// EnvKey is the environment variable holding the Kimi API key.
const EnvKey = "KIMI_API_KEY"

// Options configures the Kimi provider.
type Options struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
	Limits  provider.Limits
}

// New creates a Kimi provider.
func New(opts Options) provider.Provider {
	base := opts.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	return provider.NewOpenAICompat(provider.Options{
		Name:    "kimi",
		APIKey:  opts.APIKey,
		BaseURL: base,
		Client:  opts.Client,
		Limits:  opts.Limits,
	})
}
