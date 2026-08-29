package costtest

import (
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/cost"
)

func TestCalculateCost(t *testing.T) {
	c := cost.Calculator{}
	p := cost.Pricing{InputCostPerMillion: 1, OutputCostPerMillion: 2}
	u := cost.Usage{PromptTokens: 1_000_000, CompletionTokens: 1_000_000}
	if got := c.Calculate(u, p); got != 3 {
		t.Fatalf("expected 3, got %v", got)
	}
}

func TestManagerTrackAndTotal(t *testing.T) {
	m := cost.NewManager()
	total, err := m.TrackCost(cost.Record{Model: "a", Cost: 1.5})
	if err != nil || total != 1.5 {
		t.Fatalf("expected total 1.5, got %v err=%v", total, err)
	}
	m.TrackCost(cost.Record{Model: "a", Cost: 0.5})
	m.TrackCost(cost.Record{Model: "b", Cost: 2.0})
	if m.GetTotalCost() != 4.0 {
		t.Fatalf("expected total 4.0, got %v", m.GetTotalCost())
	}
	if m.GetModelCost("a") != 2.0 {
		t.Fatalf("expected model a cost 2.0, got %v", m.GetModelCost("a"))
	}
	m.Reset()
	if m.GetTotalCost() != 0 {
		t.Fatalf("expected reset total 0, got %v", m.GetTotalCost())
	}
}

func TestTrackerTracksUsage(t *testing.T) {
	m := cost.NewManager()
	tr := cost.NewTracker(m)
	total, err := tr.TrackUsage("opencode", "big-pickle", cost.Usage{PromptTokens: 1_000_000, CompletionTokens: 1_000_000}, cost.Pricing{InputCostPerMillion: 1, OutputCostPerMillion: 2})
	if err != nil {
		t.Fatal(err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %v", total)
	}
}

func TestThresholdAlertFiresOnce(t *testing.T) {
	var fired []string
	alert := cost.NewThresholdAlert(5.0, func(model string, total float64) {
		fired = append(fired, model)
	})
	alert.Check("a", 3.0)
	alert.Check("a", 6.0)
	alert.Check("a", 9.0)
	if len(fired) != 1 {
		t.Fatalf("expected alert to fire once, got %d", len(fired))
	}
	if e := alert.Error("a", 6.0); e == nil {
		t.Fatal("expected threshold error")
	}
}
