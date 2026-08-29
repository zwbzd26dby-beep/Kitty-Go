package llm

import (
	"context"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// MockClient is a deterministic stub that always returns a fixed response,
// ignoring history. Used for one-shot mode and as a regression baseline.
type MockClient struct{}

// NewMockClient creates a MockClient.
func NewMockClient() *MockClient { return &MockClient{} }

// Generate returns the fixed mock response.
func (m *MockClient) Generate(_ context.Context, _ types.Message, _ []types.Turn) (types.Response, error) {
	return types.NewResponse("Kitty LLM mock działa.")
}
