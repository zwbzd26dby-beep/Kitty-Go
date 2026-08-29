// Package opencode provides the OpenCode Zen provider (OpenAI-compatible).
//
// It exposes the free Zen model catalog, with "big-pickle" preferred for
// code-generation and reasoning tasks. It is a functional requirement even
// though it is not listed in the original Master Architecture document.
package opencode

import (
	"net/http"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
)

// DefaultBaseURL is the OpenCode Zen API base URL.
const DefaultBaseURL = "https://opencode.ai/zen/v1"

// Model identifiers served by OpenCode Zen.
const (
	ModelBigPickle     = "big-pickle"
	ModelMimoV2ProFree = "mimo-v2-pro-free"
	ModelMiniMaxM2Free = "minimax-m2.5-free"
)

// DefaultModel is used when no model is specified.
const DefaultModel = ModelBigPickle

// Options configures the OpenCode provider.
type Options struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
	Limits  provider.Limits
}

// New creates an OpenCode Zen provider.
func New(opts Options) provider.Provider {
	base := opts.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	return provider.NewOpenAICompat(provider.Options{
		Name:    "opencode",
		APIKey:  opts.APIKey,
		BaseURL: base,
		Client:  opts.Client,
		Limits:  opts.Limits,
	})
}
