package provider

import "net/http"

// Options for building an OpenAI-compatible provider.
type Options struct {
	Name    string
	APIKey  string
	BaseURL string
	Client  *http.Client
	Limits  Limits
}

// NewOpenAICompat builds an OpenAI-compatible *HTTPProvider pointing at the
// given base URL's /chat/completions endpoint.
func NewOpenAICompat(opts Options) *HTTPProvider {
	base := opts.BaseURL
	chat := base + "/chat/completions"
	return &HTTPProvider{
		ProviderName: opts.Name,
		APIKey:       opts.APIKey,
		BaseURL:      base,
		ChatURL:      chat,
		Client:       opts.Client,
		LimitsV:      opts.Limits,
	}
}
