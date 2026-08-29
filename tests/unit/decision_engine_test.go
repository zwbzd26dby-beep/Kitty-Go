package unit

import (
	"errors"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/decision"
)

func TestDecisionEngineCodeTask(t *testing.T) {
	e := decision.New()
	d, err := e.Analyze(decision.TaskInfo{ID: "t1", Input: "write a function in Go that fetches an API and fix the bug"})
	if err != nil {
		t.Fatal(err)
	}
	if d.Kind != "code" {
		t.Fatalf("expected code kind, got %q", d.Kind)
	}
	if !d.HasRequirement("code-gen") || !d.HasRequirement("code-reasoning") {
		t.Fatalf("missing code requirements: %+v", d.Requirements)
	}
	if !d.NeedsCapability("CodeGeneration") {
		t.Fatal("expected CodeGeneration capability")
	}
	if d.Priority == 0 {
		t.Fatal("expected non-zero priority")
	}
}

func TestDecisionEngineVisionTask(t *testing.T) {
	e := decision.New()
	d, err := e.Analyze(decision.TaskInfo{ID: "t2", Input: "describe what is in this image and the photo scene"})
	if err != nil {
		t.Fatal(err)
	}
	if d.Kind != "vision" {
		t.Fatalf("expected vision kind, got %q", d.Kind)
	}
	if !d.NeedsCapability("Vision") {
		t.Fatal("expected Vision capability")
	}
}

func TestDecisionEngineBudgetPassThrough(t *testing.T) {
	e := decision.New()
	d := e.AnalyzeOrDefault(decision.TaskInfo{ID: "t3", Input: "help me multiply 6 by 7 please", Budget: 12.5})
	if d.Kind != "math" {
		t.Fatalf("expected math kind, got %q", d.Kind)
	}
	if d.Budget != 12.5 {
		t.Fatalf("expected budget pass-through, got %f", d.Budget)
	}
}

func TestDecisionEngineDegradesToChat(t *testing.T) {
	e := decision.New()
	if _, err := e.Analyze(decision.TaskInfo{ID: "t4", Input: "hello there, how are you doing today my friend"}); !errors.Is(err, decision.ErrNoRuleMatched) {
		t.Fatalf("expected ErrNoRuleMatched, got %v", err)
	}
	d := e.AnalyzeOrDefault(decision.TaskInfo{ID: "t4", Input: "hello there, how are you doing today my friend"})
	if d.Kind != "chat" {
		t.Fatalf("expected chat fallback, got %q", d.Kind)
	}
}

func TestDecisionEngineCategoryHintOverrides(t *testing.T) {
	e := decision.New()
	d, err := e.Analyze(decision.TaskInfo{ID: "t5", Input: "write code for a sum", Category: "vision"})
	if err != nil {
		t.Fatal(err)
	}
	if d.Kind != "vision" {
		t.Fatalf("category hint should override, got %q", d.Kind)
	}
}

func TestDecisionEnginePriorityUrgent(t *testing.T) {
	e := decision.New()
	d, err := e.Analyze(decision.TaskInfo{ID: "t6", Input: "URGENT: critical production bug in this code, fix immediately"})
	if err != nil {
		t.Fatal(err)
	}
	if d.Priority < 20 {
		t.Fatalf("urgent task should score higher, got %d", d.Priority)
	}
}
