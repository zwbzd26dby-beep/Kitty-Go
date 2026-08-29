package routertest

import (
	"errors"
	"testing"

	reg "github.com/zwbzd26dby-beep/Kitty-Go/internal/model"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/router/model"
)

func buildZenRegistry(t *testing.T) *reg.Registry {
	t.Helper()
	r := reg.NewRegistry()
	if err := reg.RegisterKnown(r); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestRouterPrefersBigPickleForCode(t *testing.T) {
	r := buildZenRegistry(t)
	rt := model.NewRouter(r)
	m, err := rt.Select(model.RouteRequest{
		Task:     "write go code",
		Decision: model.Decision{RequiresCodeGeneration: true},
	})
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if m.Model != "big-pickle" {
		t.Fatalf("expected big-pickle for code, got %q", m.Model)
	}
}

func TestRouterRespectsPreferredProvider(t *testing.T) {
	r := buildZenRegistry(t)
	rt := model.NewRouter(r)
	m, err := rt.Select(model.RouteRequest{
		Decision: model.Decision{PreferredProvider: "openai"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if m.Provider != "openai" {
		t.Fatalf("expected openai, got %q", m.Provider)
	}
}

func TestRouterEmptyRegistryNoModelAvailable(t *testing.T) {
	r := reg.NewRegistry()
	rt := model.NewRouter(r)
	_, err := rt.Select(model.RouteRequest{})
	if !errors.Is(err, model.ErrNoModelAvailable) {
		t.Fatalf("expected ErrNoModelAvailable, got %v", err)
	}
}

func TestRouterBudgetFavoursCheapModels(t *testing.T) {
	r := reg.NewRegistry()
	r.Register(reg.Model{
		IDs:          reg.IDs{Provider: "demo", Model: "expensive"},
		Availability: reg.Available,
		Capabilities: reg.Capabilities{Streaming: true},
		Pricing:      reg.Pricing{InputCostPerMillion: 100, OutputCostPerMillion: 200},
	})
	r.Register(reg.Model{
		IDs:          reg.IDs{Provider: "demo", Model: "cheap"},
		Availability: reg.Available,
		Capabilities: reg.Capabilities{Streaming: true},
		Pricing:      reg.Pricing{InputCostPerMillion: 0.01, OutputCostPerMillion: 0.02},
	})
	rt := model.NewRouter(r)
	// A budget between the two costs must eliminate "expensive".
	m, err := rt.Select(model.RouteRequest{
		Decision: model.Decision{MaxCost: 0.05},
	})
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if m.Model != "cheap" {
		t.Fatalf("expected cheap model within budget, got %q", m.Model)
	}
}

func TestRouterCodeRequirementFiltersNonCodeModels(t *testing.T) {
	// A registry with no code-capable models must fail for code tasks.
	r2 := reg.NewRegistry()
	r2.Register(reg.Model{IDs: reg.IDs{Provider: "demo", Model: "no-code"}, Availability: reg.Available})
	rt := model.NewRouter(r2)
	_, err := rt.Select(model.RouteRequest{
		Decision: model.Decision{RequiresCodeGeneration: true},
	})
	if !errors.Is(err, model.ErrNoModelAvailable) {
		t.Fatalf("expected no model for code when none support it, got %v", err)
	}
}

func TestRouterFallbackExcludesModel(t *testing.T) {
	r := buildZenRegistry(t)
	rt := model.NewRouter(r)
	code := model.RouteRequest{Decision: model.Decision{RequiresCodeGeneration: true}}
	m, err := rt.GetFallback(code, "opencode/big-pickle")
	if err != nil {
		t.Fatalf("GetFallback: %v", err)
	}
	if m.Model == "big-pickle" {
		t.Fatalf("expected fallback to exclude big-pickle, got %q", m.Model)
	}
}

func TestRouterChainOrder(t *testing.T) {
	r := buildZenRegistry(t)
	rt := model.NewRouter(r)
	chain := rt.Chain(model.RouteRequest{Decision: model.Decision{RequiresCodeGeneration: true}})
	if len(chain) == 0 {
		t.Fatal("expected non-empty chain")
	}
	if chain[0].Model != "big-pickle" {
		t.Fatalf("expected big-pickle first, got %q", chain[0].Model)
	}
}
