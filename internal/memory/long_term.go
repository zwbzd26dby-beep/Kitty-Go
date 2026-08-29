package memory

import (
	"sync"
	"time"
)

// longTerm is an in-memory LongTermMemory keyed by generated ID.
type longTerm struct {
	mu      sync.RWMutex
	nextID  int
	memoirs map[string]Memory
}

// NewLongTermMemory creates a LongTermMemory.
func NewLongTermMemory() LongTermMemory {
	return &longTerm{memoirs: make(map[string]Memory)}
}

// Remember stores content and returns its ID.
func (l *longTerm) Remember(content string) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.nextID++
	id := idFromCounter(l.nextID)
	l.memoirs[id] = Memory{
		ID:        id,
		Content:   content,
		CreatedAt: time.Now(),
		Meta:      map[string]string{"kind": "long-term"},
	}
	return id, nil
}

// Recall returns a memory by ID.
func (l *longTerm) Recall(id string) (Memory, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	m, ok := l.memoirs[id]
	if !ok {
		return Memory{}, ErrNotFound
	}
	return m, nil
}

// All returns all memories.
func (l *longTerm) All() []Memory {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Memory, 0, len(l.memoirs))
	for _, m := range l.memoirs {
		out = append(out, m)
	}
	return out
}

// Forget removes a memory by ID.
func (l *longTerm) Forget(id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.memoirs[id]; !ok {
		return ErrNotFound
	}
	delete(l.memoirs, id)
	return nil
}
