package integration

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/repl"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

// buildStack assembles the full Phase 1 stack:
// Agent → Orchestrator → Executor(Local) → history-echo mock.
func buildStack() agent.Agent {
	exec := execution.NewLocalExecutor()
	orch := orchestrator.New(exec, llm.NewHistoryEchoClient())
	return agent.New(orch, agent.DefaultOptions{})
}

// TestFullFlowViaAgent verifies the layered stack end-to-end without a real
// provider: each Process call round-trips through the orchestrator.
func TestFullFlowViaAgent(t *testing.T) {
	a := buildStack()
	resp1, err := a.Process(context.Background(), "Pierwsza")
	if err != nil {
		t.Fatalf("first Process: %v", err)
	}
	if !strings.HasPrefix(resp1.Content, "[1]") {
		t.Fatalf("expected [1] on first turn, got %q", resp1.Content)
	}
	resp2, err := a.Process(context.Background(), "Druga")
	if err != nil {
		t.Fatalf("second Process: %v", err)
	}
	if !strings.HasPrefix(resp2.Content, "[3]") {
		t.Fatalf("expected [3] on second turn, got %q", resp2.Content)
	}
}

// TestFullFlowViaREPL verifies the whole stack is reachable through the REPL
// engine exposed to the user.
func TestFullFlowViaREPL(t *testing.T) {
	var out bytes.Buffer
	e := repl.NewEngineWithAgent(buildStack(), strings.NewReader("Pierwsza\nDruga\nTrzecia\n\n"), &out)
	if err := e.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"[1]", "[3]", "[5]"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("expected %q in output:\n%s", want, out.String())
		}
	}
}

// TestStackRejectsBlankInput verifies validation happens before delegation.
func TestStackRejectsBlankInput(t *testing.T) {
	a := buildStack()
	if _, err := a.Process(context.Background(), "   "); err == nil {
		t.Fatal("expected error for blank input")
	}
}
