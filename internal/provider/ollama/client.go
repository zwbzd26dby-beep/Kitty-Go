// Package ollama provides the Ollama local-model provider.
package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
)

// DefaultBaseURL is the local Ollama server address.
const DefaultBaseURL = "http://localhost:11434"

// EnvKey is the environment variable holding an optional Ollama API token.
const EnvKey = "OLLAMA_API_KEY"

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Message chatMessage `json:"message"`
	Done    bool        `json:"done"`
}

// Provider implements the local Ollama provider.
type Provider struct {
	baseURL string
	client  *http.Client
}

// Options configures the Ollama provider.
type Options struct {
	BaseURL string
	Client  *http.Client
}

// New creates an Ollama provider.
func New(opts Options) *Provider {
	base := opts.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	return &Provider{baseURL: base, client: opts.Client}
}

// Name implements provider.Provider.
func (o *Provider) Name() string { return "ollama" }

// Complete performs a chat completion against the local /api/chat endpoint.
func (o *Provider) Complete(ctx context.Context, req provider.Request) (*provider.Response, error) {
	messages := make([]chatMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		messages = append(messages, chatMessage{Role: m.Role, Content: m.Content})
	}
	body, err := json.Marshal(chatRequest{Model: req.Model, Messages: messages})
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := o.client
	if client == nil {
		client = &http.Client{}
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(httpResp.Body, 512))
		return nil, fmt.Errorf("ollama: status %d: %s", httpResp.StatusCode, string(b))
	}
	var parsed chatResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	return &provider.Response{Content: parsed.Message.Content, FinishReason: "stop"}, nil
}

// Stream is not implemented for Ollama in Phase 2.
func (o *Provider) Stream(ctx context.Context, req provider.Request) (<-chan *provider.Response, error) {
	return nil, fmt.Errorf("ollama: streaming not implemented")
}

// GetPricing implements provider.Provider (local, free).
func (o *Provider) GetPricing() provider.Pricing { return provider.Pricing{} }

// GetLimits implements provider.Provider.
func (o *Provider) GetLimits() provider.Limits { return provider.Limits{} }

// HealthCheck pings the local server root.
func (o *Provider) HealthCheck(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, o.baseURL+"/api/tags", nil)
	if err != nil {
		return err
	}
	client := o.client
	if client == nil {
		client = &http.Client{}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
