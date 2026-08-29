// Package memory implements short-term, long-term, semantic memory and a
// retrieval-augmented generation (RAG) engine (Master Architecture §19).
package memory

import "time"

// Memory is a record stored in long-term/semantic memory.
type Memory struct {
	ID        string
	Content   string
	Embedding []float64
	CreatedAt time.Time
	Meta      map[string]string
}

// Store is the common persistence contract.
type Store interface {
	Add(content string, embedding []float64) (string, error)
	Get(id string) (Memory, error)
	All() []Memory
	Delete(id string) error
}

// ShortTermMemory is the conversation-scoped working memory (Master Arch §19).
type ShortTermMemory interface {
	Add(text string) error
	Recent(n int) []string
	Clear()
}

// LongTermMemory persists memories across sessions.
type LongTermMemory interface {
	Remember(content string) (string, error)
	Recall(id string) (Memory, error)
	All() []Memory
	Forget(id string) error
}

// SemanticMemory stores and retrieves memories by vector similarity.
type SemanticMemory interface {
	Store(content string, embedding []float64) (string, error)
	// SimilaritySearch returns the top-k most similar memories.
	SimilaritySearch(query []float64, k int) []Memory
}
