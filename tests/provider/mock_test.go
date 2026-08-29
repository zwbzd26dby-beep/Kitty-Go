package providertest

import (
	"context"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
)

func TestMockProviderDeterministic(t *testing.T) {
	p := provider.NewMockProvider()
	if p.Name() != "mock" {
		t.Fatalf("expected name mock, got %q", p.Name())
	}
	resp, err := p.Complete(context.Background(), provider.Request{
		Model:    "mock-model",
		Messages: []provider.Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "Kitty LLM mock działa." {
		t.Fatalf("unexpected content: %q", resp.Content)
	}
}

func TestMockProviderHealthCheck(t *testing.T) {
	p := provider.NewMockProvider()
	if err := p.HealthCheck(context.Background()); err != nil {
		t.Fatalf("mock healthcheck should pass: %v", err)
	}
}

func TestHistoryEchoProviderReflectsMessages(t *testing.T) {
	p := provider.NewHistoryEchoProvider()
	resp, err := p.Complete(context.Background(), provider.Request{
		Messages: []provider.Message{
			{Role: "user", Content: "a"},
			{Role: "assistant", Content: "b"},
			{Role: "user", Content: "c"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(resp.Content, "[3]") {
		t.Fatalf("expected [3] prefix, got %q", resp.Content)
	}
}

func TestStreamSingleChunkMocks(t *testing.T) {
	for _, p := range []provider.Provider{provider.NewMockProvider(), provider.NewHistoryEchoProvider()} {
		ch, err := p.Stream(context.Background(), provider.Request{Messages: []provider.Message{{Role: "user", Content: "x"}}})
		if err != nil {
			t.Fatalf("%s stream: %v", p.Name(), err)
		}
		got := 0
		for range ch {
			got++
		}
		if got != 1 {
			t.Fatalf("%s stream expected 1 chunk, got %d", p.Name(), got)
		}
	}
}

func TestProviderContract(t *testing.T) {
	var _ provider.Provider = provider.NewMockProvider()
	var _ provider.Provider = provider.NewHistoryEchoProvider()
}
