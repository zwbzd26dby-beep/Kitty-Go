package memory

import (
	"fmt"
	"sync"
	"time"
)

// semantic is an in-memory SemanticMemory with brute-force cosine search.
type semantic struct {
	mu       sync.RWMutex
	nextID   int
	vectors  []Memory
	metadata map[string]string
}

// NewSemanticMemory creates a SemanticMemory.
func NewSemanticMemory() SemanticMemory {
	return &semantic{metadata: make(map[string]string)}
}

// Store adds a memory with its embedding and returns its ID.
func (s *semantic) Store(content string, embedding []float64) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	id := idFromCounter(s.nextID)
	s.vectors = append(s.vectors, Memory{
		ID:        id,
		Content:   content,
		Embedding: append([]float64(nil), embedding...),
		CreatedAt: time.Now(),
	})
	s.metadata[id] = "semantic"
	return id, nil
}

// SimilaritySearch returns the top-k memories by cosine similarity.
func (s *semantic) SimilaritySearch(query []float64, k int) []Memory {
	s.mu.RLock()
	defer s.mu.RUnlock()
	scored := make([]scoredMemory, 0, len(s.vectors))
	for _, m := range s.vectors {
		scored = append(scored, scoredMemory{mem: m, score: cosine(query, m.Embedding)})
	}
	// insertion sort descending
	for i := 1; i < len(scored); i++ {
		for j := i; j > 0 && scored[j].score > scored[j-1].score; j-- {
			scored[j], scored[j-1] = scored[j-1], scored[j]
		}
	}
	if k <= 0 || k > len(scored) {
		k = len(scored)
	}
	out := make([]Memory, 0, k)
	for i := 0; i < k; i++ {
		out = append(out, scored[i].mem)
	}
	return out
}

var _ SemanticMemory = (*semantic)(nil)

type scoredMemory struct {
	mem   Memory
	score float64
}

func idFromCounter(i int) string {
	return fmt.Sprintf("mem-%d", i)
}
