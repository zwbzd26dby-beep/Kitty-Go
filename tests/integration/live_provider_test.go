package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/provider/opencode"
)

// TestLiveOpenCodeProvider performs a real completion against OpenCode Zen.
// It is skipped in CI / without OPENCODE_API_KEY (see docs/AGENTS.md).
func TestLiveOpenCodeProvider(t *testing.T) {
	key := os.Getenv(opencode.EnvKey)
	if key == "" {
		t.Skipf("skipping live test: %s not set", opencode.EnvKey)
	}
	p := opencode.New(opencode.Options{APIKey: key})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := p.Complete(ctx, opencodeRequest())
	if err != nil {
		t.Fatalf("live completion failed: %v", err)
	}
	if resp.Content == "" {
		t.Fatal("expected non-empty content from live provider")
	}
	t.Logf("live response: %s", resp.Content)
}

func opencodeRequest() provider.Request {
	return provider.Request{
		Model:    opencode.ModelBigPickle,
		Messages: []provider.Message{{Role: "user", Content: "Respond with the single word: ok"}},
	}
}
