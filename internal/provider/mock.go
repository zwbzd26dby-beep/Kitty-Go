package provider

import (
	"context"
	"fmt"
)

// MockProvider is a deterministic offline Provider used for development and
// tests. It always returns a fixed response.
type MockProvider struct {
	// FixedContent overrides the returned content (defaults to the standard
	// Kitty mock text).
	FixedContent string
}

// NewMockProvider creates a MockProvider.
func NewMockProvider() *MockProvider {
	return &MockProvider{FixedContent: "Kitty LLM mock działa."}
}

// Name implements Provider.
func (m *MockProvider) Name() string { return "mock" }

// Complete implements Provider.
func (m *MockProvider) Complete(_ context.Context, _ Request) (*Response, error) {
	return &Response{Content: m.FixedContent}, nil
}

// Stream implements Provider by emitting a single buffered response.
func (m *MockProvider) Stream(ctx context.Context, req Request) (<-chan *Response, error) {
	resp, err := m.Complete(ctx, req)
	if err != nil {
		return nil, err
	}
	ch := make(chan *Response, 1)
	ch <- resp
	close(ch)
	return ch, nil
}

// GetPricing implements Provider.
func (m *MockProvider) GetPricing() Pricing { return Pricing{} }

// GetLimits implements Provider.
func (m *MockProvider) GetLimits() Limits { return Limits{} }

// HealthCheck implements Provider.
func (m *MockProvider) HealthCheck(context.Context) error { return nil }

// HistoryEchoProvider is a mock that reflects the length of the request
// conversation history, e.g. "[3] Kitty LLM mock działa.".
type HistoryEchoProvider struct {
	FixedSuffix string
}

// NewHistoryEchoProvider creates a HistoryEchoProvider.
func NewHistoryEchoProvider() *HistoryEchoProvider {
	return &HistoryEchoProvider{FixedSuffix: "Kitty LLM mock działa."}
}

// Name implements Provider.
func (h *HistoryEchoProvider) Name() string { return "mock-history" }

// Complete implements Provider, echoing the number of messages.
func (h *HistoryEchoProvider) Complete(_ context.Context, req Request) (*Response, error) {
	text := fmt.Sprintf("[%d] %s.", len(req.Messages), h.FixedSuffix)
	return &Response{Content: text}, nil
}

// Stream implements Provider.
func (h *HistoryEchoProvider) Stream(ctx context.Context, req Request) (<-chan *Response, error) {
	resp, err := h.Complete(ctx, req)
	if err != nil {
		return nil, err
	}
	ch := make(chan *Response, 1)
	ch <- resp
	close(ch)
	return ch, nil
}

// GetPricing implements Provider.
func (h *HistoryEchoProvider) GetPricing() Pricing { return Pricing{} }

// GetLimits implements Provider.
func (h *HistoryEchoProvider) GetLimits() Limits { return Limits{} }

// HealthCheck implements Provider.
func (h *HistoryEchoProvider) HealthCheck(context.Context) error { return nil }
