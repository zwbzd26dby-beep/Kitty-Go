package model

import (
	"errors"

	reg "github.com/zwbzd26dby-beep/Kitty-Go/internal/model"
)

// ErrNoModelAvailable is returned when routing finds no candidate matching the
// decision constraints.
var ErrNoModelAvailable = errors.New("no model available for the given constraints")

// ModelRouter selects a concrete model for a task based on a Decision
// (Master Architecture §8).
type ModelRouter interface {
	// Select returns the best matching model for req, or an error if none is
	// available within the constraints.
	Select(req RouteRequest) (Model, error)
	// GetFallback returns the best matching model excluding the given models,
	// returning an error once every candidate is exhausted.
	GetFallback(req RouteRequest, exclude ...string) (Model, error)
}

// Router implements ModelRouter against a registry.
type Router struct {
	registry ModelRegistry
}

// NewRouter builds a Router backed by reg.
func NewRouter(reg ModelRegistry) *Router {
	return &Router{registry: reg}
}

var _ ModelRouter = (*Router)(nil)

// ModelRegistry is the minimal registry surface the router needs.
type ModelRegistry interface {
	List() []reg.Model
}

// Select returns the highest-scoring candidate satisfying the decision.
func (r *Router) Select(req RouteRequest) (Model, error) {
	return r.GetFallback(req)
}

// GetFallback returns the best candidate excluding the given model keys
// (formatted as "provider/id"), progressively relaxing to cheaper/available
// alternatives.
func (r *Router) GetFallback(req RouteRequest, exclude ...string) (Model, error) {
	excluded := make(map[string]bool, len(exclude))
	for _, e := range exclude {
		excluded[e] = true
	}

	cands := r.candidates(req, excluded)
	if len(cands) == 0 {
		return Model{}, ErrNoModelAvailable
	}
	// Candidate with the highest score wins.
	best := cands[0]
	for _, c := range cands[1:] {
		if c.score > best.score {
			best = c
		}
	}
	return best.model, nil
}
