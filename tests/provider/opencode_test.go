package providertest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider/opencode"
)

// TestOpenCodeDefaultsToZenBaseURL verifies the provider uses the OpenCode
// Zen base URL by default and reports the correct name.
func TestOpenCodeDefaultsToZenBaseURL(t *testing.T) {
	p := opencode.New(opencode.Options{})
	hp := p.(*provider.HTTPProvider)
	if hp.BaseURL != opencode.DefaultBaseURL {
		t.Fatalf("expected base %q, got %q", opencode.DefaultBaseURL, hp.BaseURL)
	}
	if hp.Name() != "opencode" {
		t.Fatalf("expected name opencode, got %q", hp.Name())
	}
}

// TestOpenCodeModelCatalog verifies the free Zen model set.
func TestOpenCodeModelCatalog(t *testing.T) {
	want := map[string]bool{
		opencode.ModelBigPickle:     true,
		opencode.ModelMimoV2ProFree: true,
		opencode.ModelMiniMaxM2Free: true,
	}
	for _, m := range opencode.Models {
		if !want[m] {
			t.Fatalf("unexpected model %q in catalog", m)
		}
		want[m] = false
	}
	for m, seen := range want {
		if seen {
			t.Fatalf("model %q missing from catalog", m)
		}
	}
}

// TestOpenCodeComplete verifies a round-trip through OpenCode Zen with a
// mocked transport, sending the model in the request.
func TestOpenCodeComplete(t *testing.T) {
	srv, last := newCompatServer(t, 200, okPayload("Zen reply"))
	p := opencode.New(opencode.Options{
		APIKey:  "test-key",
		BaseURL: srv.URL,
	})
	resp, err := p.Complete(context.Background(), provider.Request{
		Model:    opencode.ModelBigPickle,
		Messages: []provider.Message{{Role: "user", Content: "hello"}},
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if resp.Content != "Zen reply" {
		t.Fatalf("unexpected content %q", resp.Content)
	}
	if last.Model != opencode.ModelBigPickle {
		t.Fatalf("expected model %q, got %q", opencode.ModelBigPickle, last.Model)
	}
}

// TestOpenCodeHealthCheck pings the configured base URL.
func TestOpenCodeHealthCheck(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}))
	t.Cleanup(srv.Close)
	p := opencode.New(opencode.Options{BaseURL: srv.URL})
	if err := p.HealthCheck(context.Background()); err != nil {
		t.Fatalf("HealthCheck: %v", err)
	}
	if hits == 0 {
		t.Fatal("expected at least one healthcheck request")
	}
}
