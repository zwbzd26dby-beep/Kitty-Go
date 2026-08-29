package llm

import (
	"context"
	"fmt"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// HistoryEchoClient is a mock that reflects the length of the conversation
// history, e.g. "[3] Kitty LLM mock działa." It is used by the REPL to make
// conversation history observable.
type HistoryEchoClient struct{}

// NewHistoryEchoClient creates a HistoryEchoClient.
func NewHistoryEchoClient() *HistoryEchoClient { return &HistoryEchoClient{} }

// Generate returns the history length followed by the mock text.
func (h *HistoryEchoClient) Generate(_ context.Context, _ types.Message, history []types.Turn) (types.Response, error) {
	text := fmt.Sprintf("[%d] Kitty LLM mock działa.", len(history))
	return types.NewResponse(text)
}
