package agent

import (
	"context"
	"fmt"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// DefaultAgent is the reference Agent implementation: it maintains an
// in-memory conversation, interprets intent and delegates execution to an
// Orchestrator.
type DefaultAgent struct {
	orchestrator orchestrator.Orchestrator
	intent       IntentAnalyzer
	conv         *types.Conversation
	hist         *historyAdapter
	contexts     ContextBuilder
}

// DefaultOptions configures a DefaultAgent.
type DefaultOptions struct {
	ID       string
	Session  string
	User     string
	Intent   IntentAnalyzer
	Contexts ContextBuilder
}

// New creates a DefaultAgent with the given orchestrator and options.
// If no Intent analyzer is provided, a default one is used.
func New(orch orchestrator.Orchestrator, opts DefaultOptions) *DefaultAgent {
	conv := types.NewConversation()
	if opts.Intent == nil {
		opts.Intent = NewIntentAnalyzer()
	}
	if opts.Contexts == nil {
		opts.Contexts = NewContextBuilder(conv, opts.ID, opts.Session, opts.User)
	}
	return &DefaultAgent{
		orchestrator: orch,
		intent:       opts.Intent,
		conv:         conv,
		hist:         &historyAdapter{conversation: conv},
		contexts:     opts.Contexts,
	}
}

// Process handles raw user input and returns the assistant response.
func (a *DefaultAgent) Process(ctx context.Context, userInput string) (*Response, error) {
	msg, err := types.NewMessage(userInput)
	if err != nil {
		return nil, err
	}
	a.hist.appendUserMessage(msg)

	intent, err := a.intent.Analyze(ctx, userInput)
	if err != nil {
		return nil, err
	}

	agentCtx, err := a.contexts.Build(userInput)
	if err != nil {
		return nil, err
	}

	task := types.Task{
		ID:      fmt.Sprintf("task-%d", a.conv.Len()),
		Input:   userInput,
		Session: agentCtx.SessionID,
		History: agentCtx.History,
	}
	result, err := a.orchestrator.Orchestrate(ctx, task)
	if err != nil {
		return nil, err
	}

	resp, err := types.NewResponse(result.Content)
	if err != nil {
		return nil, err
	}
	a.hist.appendAssistantResponse(resp)

	return &Response{
		Content:       resp.Content(),
		HistoryLen:    a.conv.Len(),
		TaskID:        task.ID,
		CurrentIntent: intent.Kind,
	}, nil
}

// GetConversation returns the underlying conversation.
func (a *DefaultAgent) GetConversation() *types.Conversation {
	return a.conv
}

// ClearHistory resets the conversation memory.
func (a *DefaultAgent) ClearHistory() {
	a.hist.clear()
}
