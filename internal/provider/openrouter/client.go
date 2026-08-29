// Package openrouter provides the OpenRouter provider (a gateway to many
// models served through a single OpenAI-compatible endpoint).
package openrouter

import (
	"net/http"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
)

// DefaultBaseURL is the OpenRouter API base URL.
const DefaultBaseURL = "https://openrouter.ai/api/v1"

// EnvKey is the environment variable holding the OpenRouter API key.
const EnvKey = "OPENROUTER_API_KEY"

// Options configures the OpenRouter provider.
type Options struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
	Limits  provider.Limits
}

// New creates an OpenRouter provider.
func New(opts Options) provider.Provider {
	base := opts.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	return provider.NewOpenAICompat(provider.Options{
		Name:    "openrouter",
		APIKey:  opts.APIKey,
		BaseURL: base,
		Client:  opts.Client,
		Limits:  opts.Limits,
	})
}
