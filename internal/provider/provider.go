// Package provider implements the provider-agnostic LLM communication layer
// (Master Architecture §9). Every external AI API (and the local mock) is
// wrapped behind the Provider interface.
package provider

import (
	"context"
	"time"
)

// Provider is the abstraction every LLM backend implements.
type Provider interface {
	// Name returns the provider identifier (e.g. "openai").
	Name() string
	// Complete performs a non-streaming completion.
	Complete(ctx context.Context, req Request) (*Response, error)
	// Stream performs a streaming completion. It may return nil, nil if
	// streaming is not supported by the implementation.
	Stream(ctx context.Context, req Request) (<-chan *Response, error)
	// GetPricing returns static pricing metadata.
	GetPricing() Pricing
	// GetLimits returns the caller-defined limits for this provider.
	GetLimits() Limits
	// HealthCheck verifies the provider is reachable/healthy.
	HealthCheck(ctx context.Context) error
}

// Request is the input to a completion call.
type Request struct {
	Model       string
	Messages    []Message
	MaxTokens   int
	Temperature float64
	Tools       []string
	Stop        []string
}

// Message is a single dialogue turn passed to a provider. Role is one of
// "system", "user" or "assistant".
type Message struct {
	Role    string
	Content string
}

// Response is the result of a completion call.
type Response struct {
	Content      string
	Usage        Usage
	FinishReason string
	Cost         float64
	Duration     time.Duration
}

// Usage reports token consumption.
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Pricing holds static price metadata for a provider.
type Pricing struct {
	InputCostPerToken  float64
	OutputCostPerToken float64
	FixedCost          float64
}

// Limits holds caller-defined limits for a provider.
type Limits struct {
	MaxRequestsPerMinute uint32
	MaxTokensPerMinute   uint32
}
