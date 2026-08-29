package unit

import (
	"context"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

// newTestAgent builds the Phase 1 layered stack with the history-echo mock.
func newTestAgent(t *testing.T) *agent.DefaultAgent {
	t.Helper()
	exec := execution.NewLocalExecutor()
	orch := orchestrator.New(exec, llm.NewHistoryEchoClient())
	return agent.New(orch, agent.DefaultOptions{})
}

func TestAgentProcessReturnsContent(t *testing.T) {
	a := newTestAgent(t)
	resp, err := a.Process(context.Background(), "Hej Kitty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Content, "mock działa") {
		t.Fatalf("unexpected content: %q", resp.Content)
	}
}

func TestAgentTracksHistory(t *testing.T) {
	a := newTestAgent(t)
	if _, err := a.Process(context.Background(), "Pierwsza"); err != nil {
		t.Fatal(err)
	}
	// After one turn the conversation has 2 turns (message + response) and
	// the response reflects history length at generation time ([1]).
	if got := a.GetConversation().Len(); got != 2 {
		t.Fatalf("expected 2 turns after one Process, got %d", got)
	}
	if _, err := a.Process(context.Background(), "Druga"); err != nil {
		t.Fatal(err)
	}
	if got := a.GetConversation().Len(); got != 4 {
		t.Fatalf("expected 4 turns after two Process calls, got %d", got)
	}
}

func TestAgentClearHistory(t *testing.T) {
	a := newTestAgent(t)
	if _, err := a.Process(context.Background(), "Pierwsza"); err != nil {
		t.Fatal(err)
	}
	a.ClearHistory()
	if got := a.GetConversation().Len(); got != 0 {
		t.Fatalf("expected empty conversation after ClearHistory, got %d", got)
	}
}

func TestAgentRejectsBlankInput(t *testing.T) {
	a := newTestAgent(t)
	for _, in := range []string{"", "   ", "\n"} {
		if _, err := a.Process(context.Background(), in); err == nil {
			t.Fatalf("expected error for %q", in)
		}
	}
}

func TestAgentResponseFields(t *testing.T) {
	a := newTestAgent(t)
	resp, err := a.Process(context.Background(), "Zadanie")
	if err != nil {
		t.Fatal(err)
	}
	if resp.TaskID == "" {
		t.Fatal("expected a non-empty TaskID")
	}
	if resp.CurrentIntent != "chat" {
		t.Fatalf("expected intent 'chat', got %q", resp.CurrentIntent)
	}
	if resp.HistoryLen != 2 {
		t.Fatalf("expected HistoryLen 2, got %d", resp.HistoryLen)
	}
}

// TestAgentSatisfiesInterface asserts DefaultAgent implements agent.Agent.
func TestAgentSatisfiesInterface(t *testing.T) {
	var _ agent.Agent = newTestAgent(t)
}
