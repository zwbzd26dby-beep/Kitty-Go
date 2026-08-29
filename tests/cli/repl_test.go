package clitest

import (
	"bytes"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/repl"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

func runREPL(input string) string {
	var out bytes.Buffer
	e := repl.NewEngine(llm.NewHistoryEchoClient(), strings.NewReader(input), &out)
	_ = e.Run()
	return out.String()
}

func TestREPLHistoryEchoGrowth(t *testing.T) {
	out := runREPL("Pierwsza\nDruga\nTrzecia\n\n")
	for _, want := range []string{"[1]", "[3]", "[5]"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestREPLEmptyLineExits(t *testing.T) {
	out := runREPL("\n")
	if !strings.Contains(out, "Hej Kitty") {
		t.Fatalf("missing greeting: %s", out)
	}
}

func TestREPLCommandsExit(t *testing.T) {
	for _, cmd := range []string{"/exit", "/quit"} {
		out := runREPL(cmd + "\n")
		if !strings.Contains(out, "Hej Kitty") {
			t.Fatalf("missing greeting for %s: %s", cmd, out)
		}
	}
}

func TestREPLHelpAndClear(t *testing.T) {
	out := runREPL("/help\n/clear\n/exit\n")
	if !strings.Contains(out, "Komendy:") {
		t.Fatalf("expected help output, got: %s", out)
	}
	if !strings.Contains(out, "wyczyszczona") {
		t.Fatalf("expected clear message, got: %s", out)
	}
}

func TestREPLErrorDoesNotKillLoop(t *testing.T) {
	var out bytes.Buffer
	e := repl.NewEngine(llm.NewHistoryEchoClient(), strings.NewReader("/unknown\nHello\n\n"), &out)
	if err := e.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "Nieznana komenda") {
		t.Fatalf("expected unknown-command message, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "[1]") {
		t.Fatalf("expected [1] after Hello, got: %s", out.String())
	}
}

func TestREPLCustomCommand(t *testing.T) {
	var out bytes.Buffer
	e := repl.NewEngine(llm.NewHistoryEchoClient(), strings.NewReader("/ping\n\n"), &out)
	e.RegisterCommand("ping", func([]string) error {
		out.WriteString("PONG\n")
		return nil
	})
	if err := e.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "PONG") {
		t.Fatalf("expected PONG, got: %s", out.String())
	}
}
