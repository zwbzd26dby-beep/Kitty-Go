package e2etest

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/budget"
	cd "github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/cost"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/limit"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	cr "github.com/zwbzd26dby-beep/Kitty-Go/internal/router/compute"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

func fullExec() execution.Executor {
	return execution.NewFullExecutor()
}

func computeDev(id string) cd.Device {
	return cd.Device{
		ID:      id,
		Address: "127.0.0.1:9100",
		Resources: cd.Resources{
			VRAMMB: 16000,
			RAMMB:  32000,
		},
		Health: cd.HealthHealthy,
	}
}

func compReg(devs ...cd.Device) *cr.Router {
	r := cd.NewRegistry()
	for _, d := range devs {
		_ = r.Register(d)
	}
	return cr.NewRouter(r)
}

func TestFullOrchestratorLocalSuccess(t *testing.T) {
	o := orchestrator.NewFull(orchestrator.FullOptions{
		Exec:   fullExec(),
		Client: llm.NewHistoryEchoClient(),
	})
	res, err := o.Orchestrate(context.Background(), types.Task{ID: "t1", Input: "Hello"})
	if err != nil {
		t.Fatalf("Orchestrate: %v", err)
	}
	if !strings.Contains(res.Content, "mock działa") {
		t.Fatalf("unexpected content %q", res.Content)
	}
}

func TestFullOrchestratorEmptyInputError(t *testing.T) {
	o := orchestrator.NewFull(orchestrator.FullOptions{Exec: fullExec(), Client: llm.NewMockClient()})
	if _, err := o.Orchestrate(context.Background(), types.Task{ID: "t2", Input: "   "}); !errors.Is(err, orchestrator.ErrEmptyTask) {
		t.Fatalf("expected ErrEmptyTask, got %v", err)
	}
}

func TestFullOrchestratorFailoverToDistributed(t *testing.T) {
	o := orchestrator.NewFull(orchestrator.FullOptions{
		Exec:       fullExec(),
		Client:     &flakyClient{fails: 3, content: "never"},
		MaxRetries: 2,
		Backoff:    1,
		Fallbacks: []execution.FallbackDevice{
			{Device: computeDev("dev-1"), Exec: &fakeRemoteExec{done: "c"}},
		},
	})
	res, err := o.Orchestrate(context.Background(), types.Task{ID: "t3", Input: "write some code please"})
	if err != nil {
		t.Fatalf("Orchestrate: %v", err)
	}
	if !strings.HasPrefix(res.Content, "remote[") {
		t.Fatalf("expected remote content, got %q", res.Content)
	}
}

func TestFullOrchestratorBudgetBlock(t *testing.T) {
	b := budget.NewManager(0.01, 0.05)
	o := orchestrator.NewFull(orchestrator.FullOptions{
		Exec:   fullExec(),
		Client: llm.NewMockClient(),
		Budget: b,
	})
	if err := b.Spend(0.01); err != nil {
		t.Fatal(err)
	}
	_, err := o.Orchestrate(context.Background(), types.Task{ID: "t4", Input: "code me something"})
	if err == nil || !strings.Contains(err.Error(), "budget") {
		t.Fatalf("expected budget error, got %v", err)
	}
}

func TestFullOrchestratorLimitBlock(t *testing.T) {
	l := limit.NewManager(limit.Limits{TokensMonthly: 100})
	o := orchestrator.NewFull(orchestrator.FullOptions{
		Exec:   fullExec(),
		Client: llm.NewMockClient(),
		Limits: l,
	})
	if err := l.CheckAndTrack(90); err != nil {
		t.Fatal(err)
	}
	_, err := o.Orchestrate(context.Background(), types.Task{ID: "t5", Input: "ask for help with math equations"})
	if err == nil || !strings.Contains(err.Error(), "rate limit") {
		t.Fatalf("expected rate limit error, got %v", err)
	}
}

func TestFullOrchestratorCostTracked(t *testing.T) {
	c := cost.NewManager()
	o := orchestrator.NewFull(orchestrator.FullOptions{
		Exec:   fullExec(),
		Client: llm.NewMockClient(),
		Cost:   c,
	})
	if _, err := o.Orchestrate(context.Background(), types.Task{ID: "t6", Input: "hello"}); err != nil {
		t.Fatal(err)
	}
	if c.GetTotalCost() <= 0 {
		t.Fatalf("expected tracked cost > 0, got %f", c.GetTotalCost())
	}
}

func TestFullOrchestratorDistributedFallbackViaComputeRouter(t *testing.T) {
	o := orchestrator.NewFull(orchestrator.FullOptions{
		Exec:        fullExec(),
		Client:      &flakyClient{fails: 99, content: "never"},
		MaxRetries:  1,
		Backoff:     1,
		Compute:     compReg(computeDev("n1")),
		Distributed: &fakeRemoteExec{done: "z"},
	})
	res, err := o.Orchestrate(context.Background(), types.Task{ID: "t7", Input: "code this up"})
	if err != nil {
		t.Fatalf("Orchestrate: %v", err)
	}
	if !strings.HasPrefix(res.Content, "remote[") {
		t.Fatalf("expected remote via compute router, got %q", res.Content)
	}
}
