package agent

import "github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"

// ContextBuilder assembles an Agent Context (conversation, session, history)
// for a given request.
type ContextBuilder interface {
	Build(userInput string) (Context, error)
}

// memoryContextBuilder builds context using an in-memory conversation.
type memoryContextBuilder struct {
	conversationID string
	sessionID      string
	userID         string
	conversation   *types.Conversation
}

// NewContextBuilder creates a ContextBuilder backed by the given conversation.
func NewContextBuilder(conversation *types.Conversation, conversationID, sessionID, userID string) ContextBuilder {
	return &memoryContextBuilder{
		conversationID: conversationID,
		sessionID:      sessionID,
		userID:         userID,
		conversation:   conversation,
	}
}

// Build assembles the current context for a user input.
func (b *memoryContextBuilder) Build(userInput string) (Context, error) {
	hist := []types.Turn{}
	if b.conversation != nil {
		hist = b.conversation.GetHistory()
	}
	return Context{
		ConversationID: b.conversationID,
		SessionID:      b.sessionID,
		UserID:         b.userID,
		History:        hist,
	}, nil
}
