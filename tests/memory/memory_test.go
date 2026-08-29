package memorytest

import (
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/memory"
)

func TestShortTermBounded(t *testing.T) {
	stm := memory.NewShortTermMemory(3)
	if err := stm.Add("a"); err != nil {
		t.Fatal(err)
	}
	if err := stm.Add("b"); err != nil {
		t.Fatal(err)
	}
	if err := stm.Add("c"); err != nil {
		t.Fatal(err)
	}
	if err := stm.Add("d"); err != nil {
		t.Fatal(err)
	}
	recent := stm.Recent(3)
	if len(recent) != 3 || recent[0] != "b" || recent[2] != "d" {
		t.Fatalf("unexpected recent: %v", recent)
	}
	stm.Clear()
	if got := stm.Recent(-1); len(got) != 0 {
		t.Fatalf("expected empty after clear, got %v", got)
	}
}

func TestLongTermCrud(t *testing.T) {
	ltm := memory.NewLongTermMemory()
	id, err := ltm.Remember("important fact")
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}
	m, err := ltm.Recall(id)
	if err != nil {
		t.Fatal(err)
	}
	if m.Content != "important fact" {
		t.Fatalf("unexpected content %q", m.Content)
	}
	if _, err := ltm.Recall("nope"); err != memory.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if len(ltm.All()) != 1 {
		t.Fatalf("expected 1 memory, got %d", len(ltm.All()))
	}
	if err := ltm.Forget("nope"); err != memory.ErrNotFound {
		t.Fatalf("expected ErrNotFound on forget, got %v", err)
	}
	if err := ltm.Forget(id); err != nil {
		t.Fatal(err)
	}
	if len(ltm.All()) != 0 {
		t.Fatalf("expected 0 memories after forget")
	}
}

func TestSemanticSimilarityRanks(t *testing.T) {
	emb := memory.NewHashEmbedder(64)
	sem := memory.NewSemanticMemory()
	docs := []string{
		"the cat sits on the mat",
		"dogs like to run in the park",
		"the quick brown fox jumps",
		"quantum physics equations",
	}
	for _, d := range docs {
		v, err := emb.Embed(d)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := sem.Store(d, v); err != nil {
			t.Fatal(err)
		}
	}
	qv, _ := emb.Embed("cat on mat")
	top := sem.SimilaritySearch(qv, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
	if !strings.Contains(top[0].Content, "cat") {
		t.Fatalf("expected cat doc first, got %q", top[0].Content)
	}
}

func TestRagEndToEnd(t *testing.T) {
	emb := memory.NewHashEmbedder(128)
	sem := memory.NewSemanticMemory()
	rag := memory.NewRAG(emb, sem, 2)
	sources := []string{
		"API rate limit is 60 requests per minute for the public tier",
		"Webhooks are delivered with HMAC-SHA256 signatures",
		"The billing cycle starts on the first of every month",
	}
	for _, s := range sources {
		if err := rag.Add(s); err != nil {
			t.Fatal(err)
		}
	}
	ctx, err := rag.Context("how many requests are allowed per minute?")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(ctx, "60") {
		t.Fatalf("expected rate limit context, got %q", ctx)
	}
}
