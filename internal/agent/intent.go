package agent

import "context"

// Intent describes the interpreted purpose of a user request.
type Intent struct {
	Kind string // e.g. "chat", "code", "tool"
}

// IntentAnalyzer classifies raw user input into an Intent.
type IntentAnalyzer interface {
	Analyze(ctx context.Context, text string) (Intent, error)
}

// defaultIntentAnalyzer is a minimal classifier used in Phase 1. It is
// replaced by the Decision Engine rules in Phase 16.
type defaultIntentAnalyzer struct{}

// NewIntentAnalyzer returns the default intent analyzer.
func NewIntentAnalyzer() IntentAnalyzer { return defaultIntentAnalyzer{} }

// Analyze returns a chat intent for all non-blank input in Phase 1.
func (defaultIntentAnalyzer) Analyze(_ context.Context, text string) (Intent, error) {
	return Intent{Kind: "chat"}, nil
}
