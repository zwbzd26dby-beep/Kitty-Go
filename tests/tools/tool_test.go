package toolstest

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/tools"
)

func TestRegistryAndAllowlist(t *testing.T) {
	reg := tools.NewRegistry()
	tools.RegisterBuiltins(reg)

	if len(reg.List()) != 3 {
		t.Fatalf("expected 3 builtins, got %d", len(reg.List()))
	}
	if !reg.IsAllowed("calculator") {
		t.Fatal("calculator should be allowed by default")
	}
	reg.Deny("calculator")
	if reg.IsAllowed("calculator") {
		t.Fatal("calculator should be denied")
	}
	if reg.IsAllowed("nope") {
		t.Fatal("unknown tool should not be allowed")
	}
	reg.Allow("calculator")
	if !reg.IsAllowed("calculator") {
		t.Fatal("calculator should be re-allowed")
	}
}

func TestExecutorRunsTool(t *testing.T) {
	reg := tools.NewRegistry()
	tools.RegisterBuiltins(reg)
	ex := tools.NewExecutor(reg)

	res, err := ex.Execute(context.Background(), "calculator", map[string]string{"expr": "2 + 3 * 4"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output != "14" {
		t.Fatalf("expected 14, got %q", res.Output)
	}
}

func TestExecutorParentheses(t *testing.T) {
	reg := tools.NewRegistry()
	reg.Register(tools.NewCalculator())
	ex := tools.NewExecutor(reg)

	res, err := ex.Execute(context.Background(), "calculator", map[string]string{"expr": "(2 + 3) * 4"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output != "20" {
		t.Fatalf("expected 20, got %q", res.Output)
	}
}

func TestExecutorDivisionByZero(t *testing.T) {
	reg := tools.NewRegistry()
	reg.Register(tools.NewCalculator())
	ex := tools.NewExecutor(reg)

	res, err := ex.Execute(context.Background(), "calculator", map[string]string{"expr": "5 / 0"})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Errored || res.Output == "" {
		t.Fatalf("expected errored result, got %+v", res)
	}
}

func TestExecutorDeniedAndMissing(t *testing.T) {
	reg := tools.NewRegistry()
	tools.RegisterBuiltins(reg)
	ex := tools.NewExecutor(reg)

	reg.Deny("sha256")
	if _, err := ex.Execute(context.Background(), "sha256", map[string]string{"input": "x"}); !errors.Is(err, tools.ErrToolDenied) {
		t.Fatalf("expected ErrToolDenied, got %v", err)
	}
	if _, err := ex.Execute(context.Background(), "ghost", nil); !errors.Is(err, tools.ErrToolNotFound) {
		t.Fatalf("expected ErrToolNotFound, got %v", err)
	}
}

func TestExecutorMissingRequiredArg(t *testing.T) {
	reg := tools.NewRegistry()
	reg.Register(tools.NewHashTool())
	ex := tools.NewExecutor(reg)
	if _, err := ex.Execute(context.Background(), "sha256", nil); err == nil {
		t.Fatal("expected missing-arg error")
	}
}

func TestExecutorSimulateMode(t *testing.T) {
	reg := tools.NewRegistry()
	reg.Register(tools.NewCalculator())
	ex := tools.NewExecutor(reg).Simulate()

	res, err := ex.Execute(context.Background(), "calculator", map[string]string{"expr": "1 / 0"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Output, "simulated") {
		t.Fatalf("expected simulated output, got %q", res.Output)
	}
}

func TestSha256Tool(t *testing.T) {
	reg := tools.NewRegistry()
	reg.Register(tools.NewHashTool())
	ex := tools.NewExecutor(reg)
	res, err := ex.Execute(context.Background(), "sha256", map[string]string{"input": "kitty"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Output) != 64 {
		t.Fatalf("SHA-256 hex should be 64 chars, got %d", len(res.Output))
	}
}

func TestResultCalledFormat(t *testing.T) {
	res := &tools.Result{Tool: "calculator", Args: map[string]string{"expr": "2+2"}, Output: "4"}
	called := res.Called()
	if !strings.Contains(called, "calculator") || !strings.Contains(called, "expr=2+2") {
		t.Fatalf("unexpected called format %q", called)
	}
}
