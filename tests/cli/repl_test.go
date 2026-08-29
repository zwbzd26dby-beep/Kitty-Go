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

func runREPLCore(input string) string {
	var out bytes.Buffer
	core := repl.NewCore(llm.NewHistoryEchoClient())
	e := repl.NewEngineWithCore(core, strings.NewReader(input), &out)
	_ = e.Run()
	return out.String()
}

func TestREPLMetricsCommand(t *testing.T) {
	out := runREPLCore("/metrics\n/exit\n")
	if !strings.Contains(out, "Brak zarejestrowanych metryk") {
		t.Fatalf("expected empty-metrics message, got: %s", out)
	}
}

func TestREPLBudgetCostEmpty(t *testing.T) {
	out := runREPLCore("/budget\n/cost\n/exit\n")
	if !strings.Contains(out, "Pozostały budżet") {
		t.Fatalf("expected budget output, got: %s", out)
	}
	if !strings.Contains(out, "Całkowity koszt") {
		t.Fatalf("expected cost output, got: %s", out)
	}
}

func TestREPLModelsListsKnown(t *testing.T) {
	out := runREPLCore("/models\n/exit\n")
	if !strings.Contains(out, "big-pickle") {
		t.Fatalf("expected big-pickle in models, got: %s", out)
	}
}

func TestREPLToolsListsBuiltins(t *testing.T) {
	out := runREPLCore("/tools\n/exit\n")
	for _, want := range []string{"calculator", "sha256", "clock"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected tool %q, got: %s", want, out)
		}
	}
}

func TestREPLModelSetCommand(t *testing.T) {
	out := runREPLCore("/model big-pickle\n/model\n/exit\n")
	if !strings.Contains(out, "Model: big-pickle") {
		t.Fatalf("expected model set, got: %s", out)
	}
}

func TestREPLMemoryEmpty(t *testing.T) {
	out := runREPLCore("/memory\n/exit\n")
	if !strings.Contains(out, "Pamięć pusta") {
		t.Fatalf("expected empty memory message, got: %s", out)
	}
}

func TestREPLDevicesEmpty(t *testing.T) {
	out := runREPLCore("/devices\n/exit\n")
	if !strings.Contains(out, "Brak urządzeń") {
		t.Fatalf("expected empty devices message, got: %s", out)
	}
}

func TestREPLSecurityStatus(t *testing.T) {
	out := runREPLCore("/security\n/exit\n")
	if !strings.Contains(out, "Sekrety") {
		t.Fatalf("expected security output, got: %s", out)
	}
}

func TestREPLHelpListsAllCommands(t *testing.T) {
	out := runREPLCore("/help\n/exit\n")
	for _, cmd := range []string{"/model", "/models", "/budget", "/cost", "/tools", "/memory", "/devices", "/metrics", "/security"} {
		if !strings.Contains(out, cmd) {
			t.Fatalf("expected %s in help, got: %s", cmd, out)
		}
	}
}
