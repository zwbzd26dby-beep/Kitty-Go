package types

// Turn is a single entry in a conversation history. It is either a Message
// (user) or a Response (assistant), discriminated by the Kind field.
type Turn struct {
	Kind  string // "message" or "response"
	value string
}

// MessageTurn builds a message turn.
func MessageTurn(msg Message) Turn { return Turn{Kind: "message", value: msg.Content()} }

// ResponseTurn builds a response turn.
func ResponseTurn(resp Response) Turn { return Turn{Kind: "response", value: resp.Content()} }

// IsMessage reports whether this turn is a user message.
func (t Turn) IsMessage() bool { return t.Kind == "message" }

// Content returns the underlying content.
func (t Turn) Content() string { return t.value }
