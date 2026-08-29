package types

// Conversation stores an ordered history of messages and responses in memory.
type Conversation struct {
	history []Turn
}

// NewConversation creates an empty conversation.
func NewConversation() *Conversation {
	return &Conversation{history: []Turn{}}
}

// AddMessage appends a user message to the history.
func (c *Conversation) AddMessage(msg Message) {
	c.history = append(c.history, MessageTurn(msg))
}

// AddResponse appends an assistant response to the history.
func (c *Conversation) AddResponse(resp Response) {
	c.history = append(c.history, ResponseTurn(resp))
}

// GetHistory returns a copy of the history so callers cannot mutate internals.
func (c *Conversation) GetHistory() []Turn {
	out := make([]Turn, len(c.history))
	copy(out, c.history)
	return out
}

// LastMessage returns the most recent message turn, or zero value if none.
func (c *Conversation) LastMessage() (Turn, bool) {
	for i := len(c.history) - 1; i >= 0; i-- {
		if c.history[i].IsMessage() {
			return c.history[i], true
		}
	}
	return Turn{}, false
}

// LastResponse returns the most recent response turn, or zero value if none.
func (c *Conversation) LastResponse() (Turn, bool) {
	for i := len(c.history) - 1; i >= 0; i-- {
		if c.history[i].Kind == "response" {
			return c.history[i], true
		}
	}
	return Turn{}, false
}

// Clear removes all history.
func (c *Conversation) Clear() {
	c.history = c.history[:0]
}

// Len returns the number of turns in the conversation.
func (c *Conversation) Len() int { return len(c.history) }
