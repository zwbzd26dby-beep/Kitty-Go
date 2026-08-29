package unit

import (
	"context"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// TestMockClientDeterministic verifies MockClient always returns the fixed
// response regardless of input/history.
func TestMockClientDeterministic(t *testing.T) {
	client := llm.NewMockClient()
	msg, _ := types.NewMessage("cokolwiek")
	hist := []types.Turn{types.MessageTurn(msg)}
	resp, err := client.Generate(context.Background(), msg, hist)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content() != "Kitty LLM mock działa." {
		t.Fatalf("unexpected mock output: %q", resp.Content())
	}
}

func TestHistoryEchoReflectsLength(t *testing.T) {
	client := llm.NewHistoryEchoClient()
	msg, _ := types.NewMessage("m")
	hist := []types.Turn{
		types.MessageTurn(mustMessage(t, "a")),
		types.ResponseTurn(mustResponse(t, "b")),
		types.MessageTurn(mustMessage(t, "c")),
	}
	resp, err := client.Generate(context.Background(), msg, hist)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(resp.Content(), "[3]") {
		t.Fatalf("expected [3] prefix, got %q", resp.Content())
	}
}

func TestHistoryEchoWithEmptyHistory(t *testing.T) {
	client := llm.NewHistoryEchoClient()
	msg, _ := types.NewMessage("m")
	resp, _ := client.Generate(context.Background(), msg, nil)
	if !strings.HasPrefix(resp.Content(), "[0]") {
		t.Fatalf("expected [0] prefix, got %q", resp.Content())
	}
}

// TestClientContract ensures both mocks satisfy llm.Client.
func TestClientContract(t *testing.T) {
	var _ llm.Client = llm.NewMockClient()
	var _ llm.Client = llm.NewHistoryEchoClient()
}

func mustResponse(t *testing.T, s string) types.Response {
	t.Helper()
	r, err := types.NewResponse(s)
	if err != nil {
		t.Fatalf("NewResponse(%q): %v", s, err)
	}
	return r
}
