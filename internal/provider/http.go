package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/utils"
)

// DefaultHTTPClient returns a shared http.Client with sensible timeouts.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
	}
}

// HTTPProvider is a base for OpenAI-compatible (chat/completions) providers.
// It carries a pluggable *http.Client so tests can substitute a mock transport.
type HTTPProvider struct {
	ProviderName string
	APIKey       string
	BaseURL      string
	ChatURL      string
	Client       *http.Client

	LimitsV Limits
}

func (h *HTTPProvider) Name() string { return h.ProviderName }

// chatRequest mirrors the OpenAI chat/completions request shape.
type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	Temperature *float64      `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatResponse mirrors the OpenAI chat/completions response shape.
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Complete performs a non-streaming completion against the chat/completions
// endpoint. It retries transient failures with exponential backoff and honors
// the context deadline.
func (h *HTTPProvider) Complete(ctx context.Context, req Request) (*Response, error) {
	var resp *Response
	err := utils.Retry(ctx, 3, 200*time.Millisecond, func() error {
		var derr error
		resp, derr = h.completeOnce(ctx, req)
		return derr
	})
	return resp, err
}

func (h *HTTPProvider) completeOnce(ctx context.Context, req Request) (*Response, error) {
	body, err := h.buildBody(req, false)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, h.ChatURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if h.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+h.APIKey)
	}

	httpResp, err := h.do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(httpResp.Body, 512))
		return nil, fmt.Errorf("provider %s: status %d: %s", h.Name(), httpResp.StatusCode, string(b))
	}

	var parsed chatResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("provider %s: decode: %w", h.Name(), err)
	}
	return &Response{
		Content:      parsed.Choices[0].Message.Content,
		FinishReason: parsed.Choices[0].FinishReason,
		Usage: Usage{
			PromptTokens:     parsed.Usage.PromptTokens,
			CompletionTokens: parsed.Usage.CompletionTokens,
			TotalTokens:      parsed.Usage.TotalTokens,
		},
	}, nil
}

// Stream is not implemented by the base provider (Phase 3+).
func (h *HTTPProvider) Stream(ctx context.Context, req Request) (<-chan *Response, error) {
	return nil, fmt.Errorf("provider %s: streaming not implemented", h.Name())
}

// GetPricing implements Provider.
func (h *HTTPProvider) GetPricing() Pricing { return Pricing{} }

// GetLimits implements Provider.
func (h *HTTPProvider) GetLimits() Limits { return h.LimitsV }

// HealthCheck performs a HEAD/GET to the base URL.
func (h *HTTPProvider) HealthCheck(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, h.BaseURL, nil)
	if err != nil {
		return err
	}
	resp, err := h.do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// do executes the request through the (possibly mocked) client.
func (h *HTTPProvider) do(req *http.Request) (*http.Response, error) {
	if h.Client == nil {
		h.Client = DefaultHTTPClient()
	}
	return h.Client.Do(req)
}

func (h *HTTPProvider) buildBody(req Request, stream bool) ([]byte, error) {
	messages := make([]chatMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		messages = append(messages, chatMessage{Role: m.Role, Content: m.Content})
	}
	maxTokens := req.MaxTokens
	temperature := req.Temperature
	payload := chatRequest{
		Model:       req.Model,
		Messages:    messages,
		MaxTokens:   &maxTokens,
		Temperature: &temperature,
		Stream:      stream,
	}
	if req.MaxTokens == 0 {
		payload.MaxTokens = nil
	}
	if req.Temperature == 0 {
		payload.Temperature = nil
	}
	return json.Marshal(payload)
}
