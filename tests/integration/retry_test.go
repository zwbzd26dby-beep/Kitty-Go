package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider/openai"
)

// TestProviderRetriesOnServerError verifies that a 5xx response is retried by
// HTTPProvider.Complete and that the second attempt succeeds.
func TestProviderRetriesOnServerError(t *testing.T) {
	var mu sync.Mutex
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		attempts++
		n := attempts
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if n < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]interface{}{"content": "recovered"}, "finish_reason": "stop"},
			},
		})
	}))
	t.Cleanup(srv.Close)

	p := openai.New(openai.Options{BaseURL: srv.URL})
	resp, err := p.Complete(context.Background(), provider.Request{
		Messages: []provider.Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("expected retry to succeed: %v", err)
	}
	if resp.Content != "recovered" {
		t.Fatalf("unexpected content %q", resp.Content)
	}
	mu.Lock()
	defer mu.Unlock()
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

// TestProviderDoesNotRetryClientError verifies a 4xx is not retried.
func TestProviderDoesNotRetryClientError(t *testing.T) {
	var mu sync.Mutex
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		attempts++
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	}))
	t.Cleanup(srv.Close)

	p := openai.New(openai.Options{BaseURL: srv.URL})
	_, err := p.Complete(context.Background(), provider.Request{
		Messages: []provider.Message{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	mu.Lock()
	defer mu.Unlock()
	if attempts != 1 {
		t.Fatalf("expected 1 attempt for 4xx, got %d", attempts)
	}
}
