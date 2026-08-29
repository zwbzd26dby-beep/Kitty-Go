// Package model implements the Model Router (Master Architecture §8) which
// selects a concrete model from the registry based on a Decision.
package model

import reg "github.com/zwbzd26dby-beep/Kitty-Go/internal/model"

// Model is an alias for a registry model entry.
type Model = reg.Model

// Candidate is a scored model during selection.
type candidate struct {
	model reg.Model
	score float64
}

// Decision captures the routing requirements produced by the Decision Engine
// (Master Architecture §7, implemented in Phase 16). It is defined here so the
// router can be used standalone; Phase 16 populates it.
type Decision struct {
	// RequiresCodeGeneration flags code-generation/reasoning tasks, which in
	// this deployment prefer OpenCode's big-pickle model.
	RequiresCodeGeneration bool
	// Priority weights the selection: "cost", "latency" or "quality".
	Priority string
	// MaxCost caps the per-request spend (0 = unlimited).
	MaxCost float64
	// PreferredProvider, when set, restricts candidates to that provider.
	PreferredProvider string
	// PreferredModel, when set, is tried first before scoring.
	PreferredModel string
}

// RouteRequest bundles a task with its Decision for routing.
type RouteRequest struct {
	Task     string
	Decision Decision
}
