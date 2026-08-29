// Package agent implements the Agent Core: it understands user intent,
// manages conversation context and history, and hands tasks to the
// Orchestrator.
package agent

import (
	"context"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// Agent is the central entry point for a user request.
type Agent interface {
	// Process handles raw user input and returns the assistant response.
	Process(ctx context.Context, userInput string) (*Response, error)
	// GetConversation returns the underlying conversation (may be nil).
	GetConversation() *types.Conversation
	// ClearHistory resets the conversation memory.
	ClearHistory()
}

// Response is the result of processing a user request.
type Response struct {
	Content       string
	HistoryLen    int
	TaskID        string
	CurrentIntent string
}

// Context describes the current session context for a request.
type Context struct {
	ConversationID string
	SessionID      string
	UserID         string
	History        []types.Turn
	CurrentIntent  string
}
