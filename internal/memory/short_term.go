package memory

import (
	"sync"
)

// shortTerm is an in-memory ShortTermMemory (bounded ring).
type shortTerm struct {
	mu   sync.RWMutex
	cap  int
	text []string
}

// NewShortTermMemory creates a bounded ShortTermMemory.
func NewShortTermMemory(cap int) ShortTermMemory {
	return &shortTerm{cap: cap}
}

// Add appends text, trimming the oldest when over capacity.
func (s *shortTerm) Add(text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.text) >= s.cap {
		s.text = s.text[1:]
	}
	s.text = append(s.text, text)
	return nil
}

// Recent returns the last n entries (all if n <= 0 or n > len).
func (s *shortTerm) Recent(n int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if n <= 0 || n > len(s.text) {
		n = len(s.text)
	}
	out := make([]string, n)
	copy(out, s.text[len(s.text)-n:])
	return out
}

// Clear empties the memory.
func (s *shortTerm) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.text = s.text[:0]
}
