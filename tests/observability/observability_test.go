package obsertest

import (
	"bytes"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/observability"
)

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	l := observability.NewLogger(&buf, observability.LevelInfo)
	l.Debug("hidden", "k", "v")
	l.Info("shown")
	l.Error("bad", "err", "x")
	out := buf.String()
	if strings.Contains(out, "hidden") {
		t.Fatalf("debug should be filtered: %q", out)
	}
	if !strings.Contains(out, "shown") || !strings.Contains(out, "bad") {
		t.Fatalf("expected info/error lines, got %q", out)
	}
}

func TestMetricsCounterAndGauge(t *testing.T) {
	m := observability.NewMetrics()
	m.Inc("requests")
	m.Inc("requests")
	m.Add("tokens", 5)
	m.Set("cpu", 0.42)
	if m.Get("requests") != 2 {
		t.Fatalf("expected 2 requests, got %d", m.Get("requests"))
	}
	cs, gs := m.Snapshot()
	if cs["tokens"] != 5 || gs["cpu"] != 0.42 {
		t.Fatalf("snapshot mismatch: %v %v", cs, gs)
	}
}

func TestUsageAggregation(t *testing.T) {
	u := observability.NewUsage()
	u.LogP("big-pickle", 100, 50)
	u.LogP("big-pickle", 200, 70)
	u.LogP("other", 10, 10)
	total := u.Total()
	if total["big-pickle"].PromptTokens != 300 || total["big-pickle"].CompletionTokens != 120 {
		t.Fatalf("unexpected totals: %+v", total["big-pickle"])
	}
	if total["big-pickle"].Calls != 2 {
		t.Fatalf("expected 2 calls, got %d", total["big-pickle"].Calls)
	}
	if len(u.Records()) != 3 {
		t.Fatalf("expected 3 records, got %d", len(u.Records()))
	}
}

func TestCostComputation(t *testing.T) {
	c := observability.NewCost(map[string]observability.PricePoint{
		"m": {PromptIn: 0.15, CompletionIn: 0.60},
	})
	c.LogP("m", 1_000_000, 1_000_000)
	total := c.TotalSpend()
	// 0.15 + 0.60 = 0.75
	if total["m"] < 0.749 || total["m"] > 0.751 {
		t.Fatalf("expected ~0.75, got %f", total["m"])
	}
	if len(c.Records()) != 1 {
		t.Fatal("expected 1 record")
	}
}

func TestTraceSpans(t *testing.T) {
	tr := observability.StartTrace()
	tr.Start("route")
	tr.Start("execute")
	tr.End()
	tr.End()
	spans := tr.Drain()
	if len(spans) != 2 {
		t.Fatalf("expected 2 spans, got %d", len(spans))
	}
	if spans[0].Name != "execute" || spans[1].Name != "route" {
		t.Fatalf("expected LIFO span order, got %+v", spans)
	}
}

func TestAuditSinkAndFilter(t *testing.T) {
	a := observability.NewAudit()
	var got []string
	a.OnEvent(func(ev observability.AuditEvent) { got = append(got, ev.Action) })
	a.Record("alice", "login", "session start")
	a.Record("alice", "logout", "session end")
	a.Record("bob", "login", "another")
	if len(got) != 3 {
		t.Fatalf("sink should see 3 events, got %v", got)
	}
	if len(a.Filter("login")) != 2 {
		t.Fatalf("expected 2 logins")
	}
	if len(a.Filter("logout")) != 1 {
		t.Fatalf("expected 1 logout")
	}
}

func TestObservabilityRoot(t *testing.T) {
	var buf bytes.Buffer
	o := observability.New(&buf)
	o.Metrics.Inc("calls")
	if o.Metrics.Get("calls") != 1 {
		t.Fatal("metrics root broken")
	}
	o.Audit.Record("user", "action", "d")
	if len(o.Audit.Events()) != 1 {
		t.Fatal("audit root broken")
	}
	o.Start("work")
	o.End()
	if len(o.Trace.Drain()) != 1 {
		t.Fatal("trace root broken")
	}
}
