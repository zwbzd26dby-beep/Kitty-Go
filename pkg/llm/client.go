// Package llm defines the provider-agnostic LLM client contract and the
// deterministic mocks used for offline development and testing.
package llm

import (
	"context"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// Client is the abstraction every LLM provider / mock implements.
type Client interface {
	// Generate produces a response for the given message, optionally
	// considering the conversation history.
	Generate(ctx context.Context, msg types.Message, history []types.Turn) (types.Response, error)
}
