package providertest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider/ollama"
)

// TestOllamaComplete verifies chat completion against the local /api/chat.
func TestOllamaComplete(t *testing.T) {
	var gotModel string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		var body struct {
			Model string `json:"model"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		gotModel = body.Model
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": map[string]interface{}{"role": "assistant", "content": "Ollama reply"},
			"done":    true,
		})
	}))
	t.Cleanup(srv.Close)

	p := ollama.New(ollama.Options{BaseURL: srv.URL})
	resp, err := p.Complete(context.Background(), provider.Request{
		Model:    ollama.DefaultModel,
		Messages: []provider.Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if resp.Content != "Ollama reply" {
		t.Fatalf("unexpected content %q", resp.Content)
	}
	if gotModel != ollama.DefaultModel {
		t.Fatalf("expected model %q, got %q", ollama.DefaultModel, gotModel)
	}
}

func TestOllamaName(t *testing.T) {
	if got := ollama.New(ollama.Options{}).Name(); got != "ollama" {
		t.Fatalf("expected ollama, got %q", got)
	}
}

func TestOllamaHealthCheck(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})
	}))
	t.Cleanup(srv.Close)
	p := ollama.New(ollama.Options{BaseURL: srv.URL})
	if err := p.HealthCheck(context.Background()); err != nil {
		t.Fatalf("HealthCheck: %v", err)
	}
	if hits == 0 {
		t.Fatal("expected healthcheck request")
	}
}
