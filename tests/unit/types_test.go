package unit

import (
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

func TestNewMessage(t *testing.T) {
	msg, err := types.NewMessage("Hej Kitty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Content() != "Hej Kitty" {
		t.Fatalf("content mismatch: %q", msg.Content())
	}
}

func TestNewMessageRejectsEmpty(t *testing.T) {
	for _, in := range []string{"", "   ", "\n", "\t"} {
		if _, err := types.NewMessage(in); err == nil {
			t.Fatalf("expected error for %q", in)
		}
	}
}

func TestNewResponseValidatesContent(t *testing.T) {
	if _, err := types.NewResponse("ok"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := types.NewResponse("   "); err == nil {
		t.Fatal("expected error for blank response")
	}
}

func TestMessageResponseImmutable(t *testing.T) {
	msg, _ := types.NewMessage("abc")
	resp, _ := types.NewResponse("xyz")
	// Only accessors are exported; content cannot be mutated from outside.
	if msg.Content() != "abc" || resp.Content() != "xyz" {
		t.Fatal("content changed unexpectedly")
	}
}

func TestConversationHistoryDefensiveCopy(t *testing.T) {
	c := types.NewConversation()
	msg, _ := types.NewMessage("Pierwsza")
	resp, _ := types.NewResponse("[1] Kitty LLM mock działa.")
	c.AddMessage(msg)
	c.AddResponse(resp)
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
	hist := c.GetHistory()
	hist[0] = types.MessageTurn(mustMessage(t, "mutated"))
	if got := c.GetHistory()[0]; got.Content() == "mutated" {
		t.Fatal("history is not a defensive copy")
	}
}

func TestConversationLastAndClear(t *testing.T) {
	c := types.NewConversation()
	c.AddMessage(mustMessage(t, "m1"))
	c.AddMessage(mustMessage(t, "m2"))
	last, ok := c.LastMessage()
	if !ok || last.Content() != "m2" {
		t.Fatalf("expected m2, got %q ok=%v", last.Content(), ok)
	}
	if _, ok := c.LastResponse(); ok {
		t.Fatal("expected no response yet")
	}
	r, _ := types.NewResponse("r1")
	c.AddResponse(r)
	if last, ok := c.LastResponse(); !ok || last.Content() != "r1" {
		t.Fatalf("expected r1, got %q", last.Content())
	}
	c.Clear()
	if c.Len() != 0 {
		t.Fatalf("expected empty after clear, got %d", c.Len())
	}
}

func TestTaskCarriesFields(t *testing.T) {
	tsk := types.Task{ID: "t1", Input: "input", Session: "s1"}
	if tsk.ID != "t1" || tsk.Input != "input" || tsk.Session != "s1" {
		t.Fatal("task fields mismatch")
	}
}

func mustMessage(t *testing.T, s string) types.Message {
	t.Helper()
	m, err := types.NewMessage(s)
	if err != nil {
		t.Fatalf("NewMessage(%q): %v", s, err)
	}
	return m
}
