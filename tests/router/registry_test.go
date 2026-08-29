package routertest

import (
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/model"
)

func TestRegistryRegisterGet(t *testing.T) {
	r := model.NewRegistry()
	m := model.Model{
		IDs:           model.IDs{Provider: "opencode", Model: "big-pickle"},
		ContextWindow: 128000,
		Availability:  model.Available,
	}
	if err := r.Register(m); err != nil {
		t.Fatal(err)
	}
	got, err := r.Get("opencode", "big-pickle")
	if err != nil {
		t.Fatal(err)
	}
	if got.ContextWindow != 128000 {
		t.Fatalf("unexpected context window %d", got.ContextWindow)
	}
}

func TestRegistryRegisterRequiresIDs(t *testing.T) {
	r := model.NewRegistry()
	if err := r.Register(model.Model{}); err == nil {
		t.Fatal("expected error for empty model")
	}
}

func TestRegistryGetMissing(t *testing.T) {
	r := model.NewRegistry()
	if _, err := r.Get("openai", "nope"); err == nil {
		t.Fatal("expected error for missing model")
	}
}

func TestRegistryListAndFilter(t *testing.T) {
	r := model.NewRegistry()
	r.Register(model.Model{IDs: model.IDs{Provider: "opencode", Model: "a"}, Availability: model.Available})
	r.Register(model.Model{IDs: model.IDs{Provider: "opencode", Model: "b"}, Availability: model.Available})
	r.Register(model.Model{IDs: model.IDs{Provider: "openai", Model: "c"}, Availability: model.Available})
	if got := len(r.List()); got != 3 {
		t.Fatalf("expected 3 models, got %d", got)
	}
	if got := len(r.ListByProvider("opencode")); got != 2 {
		t.Fatalf("expected 2 opencode models, got %d", got)
	}
}

func TestRegistryUpdatePricingAndAvailability(t *testing.T) {
	r := model.NewRegistry()
	r.Register(model.Model{IDs: model.IDs{Provider: "opencode", Model: "big-pickle"}, Availability: model.Available})
	if err := r.UpdatePricing("opencode", "big-pickle", model.Pricing{InputCostPerMillion: 1, OutputCostPerMillion: 2}); err != nil {
		t.Fatal(err)
	}
	if err := r.UpdateAvailability("opencode", "big-pickle", model.Unavailable); err != nil {
		t.Fatal(err)
	}
	m, _ := r.Get("opencode", "big-pickle")
	if m.Pricing.InputCostPerMillion != 1 {
		t.Fatalf("pricing not updated: %+v", m.Pricing)
	}
	if m.Availability != model.Unavailable {
		t.Fatalf("availability not updated: %q", m.Availability)
	}
}

func TestRegisterKnownSeedsZenCatalog(t *testing.T) {
	r := model.NewRegistry()
	if err := model.RegisterKnown(r); err != nil {
		t.Fatal(err)
	}
	zen := r.ListByProvider("opencode")
	if len(zen) != 3 {
		t.Fatalf("expected 3 opencode zen models, got %d", len(zen))
	}
	if _, err := r.Get("opencode", "big-pickle"); err != nil {
		t.Fatalf("big-pickle not registered: %v", err)
	}
}

func TestPricingCostFor(t *testing.T) {
	p := model.Pricing{InputCostPerMillion: 1, OutputCostPerMillion: 2}
	// 1M in, 1M out => 1*1 + 1*2 = 3
	if got := p.CostFor(1_000_000, 1_000_000); got != 3 {
		t.Fatalf("expected 3, got %v", got)
	}
}
