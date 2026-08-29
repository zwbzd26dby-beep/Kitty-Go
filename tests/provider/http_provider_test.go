package providertest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider/deepseek"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider/kimi"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider/openai"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider/openrouter"
)

// newCompatServer returns an httptest server that answers OpenAI-style
// chat/completions and records the last request.
func newCompatServer(t *testing.T, status int, payload map[string]interface{}) (*httptest.Server, *chatCompatRequest) {
	t.Helper()
	var lastRequest chatCompatRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&lastRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &lastRequest
}

type chatCompatRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

func okPayload(content string) map[string]interface{} {
	return map[string]interface{}{
		"choices": []map[string]interface{}{
			{"message": map[string]interface{}{"content": content}, "finish_reason": "stop"},
		},
		"usage": map[string]interface{}{
			"prompt_tokens":     5,
			"completion_tokens": 3,
			"total_tokens":      8,
		},
	}
}

func TestOpenAICompatProviders(t *testing.T) {
	providers := map[string]provider.Provider{
		"openai":     openai.New(openai.Options{}),
		"kimi":       kimi.New(kimi.Options{}),
		"deepseek":   deepseek.New(deepseek.Options{}),
		"openrouter": openrouter.New(openrouter.Options{}),
	}
	for name := range providers {
		srv, last := newCompatServer(t, 200, okPayload("Hello from "+name))
		p := providers[name]
		// Rebuild the provider pointing at the test server.
		p = rebuildCompat(name, srv.URL)
		resp, err := p.Complete(context.Background(), provider.Request{
			Model:    "some-model",
			Messages: []provider.Message{{Role: "user", Content: "hi"}},
		})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if resp.Content != "Hello from "+name {
			t.Fatalf("%s: unexpected content %q", name, resp.Content)
		}
		if resp.Usage.TotalTokens != 8 {
			t.Fatalf("%s: unexpected usage %+v", name, resp.Usage)
		}
		if last.Model != "some-model" {
			t.Fatalf("%s: request model %q", name, last.Model)
		}
	}
}

// TestCompatProviderSendsAuthHeader verifies a Bearer token is attached.
func TestCompatProviderSendsAuthHeader(t *testing.T) {
	var authHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(okPayload("ok"))
	}))
	t.Cleanup(srv.Close)

	p := openai.New(openai.Options{APIKey: "secret-123", BaseURL: srv.URL})
	if _, err := p.Complete(context.Background(), provider.Request{Messages: []provider.Message{{Role: "user", Content: "hi"}}}); err != nil {
		t.Fatal(err)
	}
	if authHeader != "Bearer secret-123" {
		t.Fatalf("expected Bearer header, got %q", authHeader)
	}
}

// TestCompatProviderErrorStatus verifies non-2xx responses surface errors.
func TestCompatProviderErrorStatus(t *testing.T) {
	srv, _ := newCompatServer(t, 401, nil)
	p := openai.New(openai.Options{BaseURL: srv.URL})
	_, err := p.Complete(context.Background(), provider.Request{Messages: []provider.Message{{Role: "user", Content: "hi"}}})
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Fatalf("expected 401 in error, got %v", err)
	}
}

func rebuildCompat(name, baseURL string) provider.Provider {
	switch name {
	case "openai":
		return openai.New(openai.Options{BaseURL: baseURL})
	case "kimi":
		return kimi.New(kimi.Options{BaseURL: baseURL})
	case "deepseek":
		return deepseek.New(deepseek.Options{BaseURL: baseURL})
	case "openrouter":
		return openrouter.New(openrouter.Options{BaseURL: baseURL})
	}
	return nil
}
