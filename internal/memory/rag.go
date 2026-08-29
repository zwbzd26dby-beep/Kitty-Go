package memory

import (
	"errors"
	"strings"
)

// ErrNotFound is returned when a memory ID does not exist.
var ErrNotFound = errors.New("memory not found")

// RAG retrieves relevant context for a query by embedding + similarity search.
type RAG struct {
	emb  Embedder
	mem  SemanticMemory
	topK int
}

// NewRAG builds a RAG engine over the given semantic store.
func NewRAG(emb Embedder, mem SemanticMemory, topK int) *RAG {
	if topK <= 0 {
		topK = 3
	}
	return &RAG{emb: emb, mem: mem, topK: topK}
}

// Query embeds the question and returns the top-k relevant memory contents.
func (r *RAG) Query(query string) ([]string, error) {
	vec, err := r.emb.Embed(query)
	if err != nil {
		return nil, err
	}
	results := r.mem.SimilaritySearch(vec, r.topK)
	out := make([]string, 0, len(results))
	for _, m := range results {
		out = append(out, m.Content)
	}
	return out, nil
}

// Context folds the top-k results into a single context string.
func (r *RAG) Context(query string) (string, error) {
	docs, err := r.Query(query)
	if err != nil {
		return "", err
	}
	return strings.Join(docs, "\n---\n"), nil
}

// Add stores a document in the semantic store.
func (r *RAG) Add(content string) error {
	vec, err := r.emb.Embed(content)
	if err != nil {
		return err
	}
	_, err = r.mem.Store(content, vec)
	return err
}
