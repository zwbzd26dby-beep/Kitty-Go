package model

import (
	"fmt"
	"sync"
)

// ModelRegistry is the central store of models (Master Architecture §10).
type ModelRegistry interface {
	Register(m Model) error
	Get(provider, id string) (Model, error)
	List() []Model
	ListByProvider(provider string) []Model
	UpdatePricing(provider, id string, p Pricing) error
	UpdateAvailability(provider, id string, a Availability) error
}

// Registry is a thread-safe in-memory ModelRegistry.
type Registry struct {
	mu     sync.RWMutex
	models map[string]Model
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{models: make(map[string]Model)}
}

func key(provider, id string) string { return provider + "/" + id }

// Register adds or replaces a model.
func (r *Registry) Register(m Model) error {
	if m.Provider == "" || m.Model == "" {
		return fmt.Errorf("model provider and id are required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.models[key(m.Provider, m.Model)] = m
	return nil
}

// Get returns a model by provider and id.
func (r *Registry) Get(provider, id string) (Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.models[key(provider, id)]
	if !ok {
		return Model{}, fmt.Errorf("model %q not found in provider %q", id, provider)
	}
	return m, nil
}

// List returns all models.
func (r *Registry) List() []Model {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Model, 0, len(r.models))
	for _, m := range r.models {
		out = append(out, m)
	}
	return out
}

// ListByProvider returns all models for a provider.
func (r *Registry) ListByProvider(provider string) []Model {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []Model
	for _, m := range r.models {
		if m.Provider == provider {
			out = append(out, m)
		}
	}
	return out
}

// UpdatePricing changes a model's pricing.
func (r *Registry) UpdatePricing(provider, id string, p Pricing) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.models[key(provider, id)]
	if !ok {
		return fmt.Errorf("model %q/%q not found", provider, id)
	}
	m.Pricing = p
	r.models[key(provider, id)] = m
	return nil
}

// UpdateAvailability changes a model's availability.
func (r *Registry) UpdateAvailability(provider, id string, a Availability) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.models[key(provider, id)]
	if !ok {
		return fmt.Errorf("model %q/%q not found", provider, id)
	}
	m.Availability = a
	r.models[key(provider, id)] = m
	return nil
}

var _ ModelRegistry = (*Registry)(nil)
