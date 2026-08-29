// Package openai provides the OpenAI ChatGPT provider.
package openai

import (
	"net/http"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
)

// DefaultBaseURL is the OpenAI API base URL.
const DefaultBaseURL = "https://api.openai.com/v1"

// Options configures the OpenAI provider.
type Options struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
	Limits  provider.Limits
}

// New creates an OpenAI provider.
func New(opts Options) provider.Provider {
	base := opts.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	return provider.NewOpenAICompat(provider.Options{
		Name:    "openai",
		APIKey:  opts.APIKey,
		BaseURL: base,
		Client:  opts.Client,
		Limits:  opts.Limits,
	})
}
