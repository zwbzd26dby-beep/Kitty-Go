package model

import (
	reg "github.com/zwbzd26dby-beep/Kitty-Go/internal/model"
)

// BigPickleModel is the OpenCode Zen model preferred for coding tasks.
const BigPickleModel = "big-pickle"

// candidates filters the registry using the decision constraints, then scores
// the survivors and returns them sorted by descending score.
func (r *Router) candidates(req RouteRequest, excluded map[string]bool) []candidate {
	var cands []candidate
	for _, m := range r.registry.List() {
		key := m.Provider + "/" + m.Model
		if excluded[key] {
			continue
		}
		if !filter(m, req.Decision) {
			continue
		}
		c := candidate{model: m, score: score(m, req.Decision)}
		cands = append(cands, c)
	}
	// Stable sort: higher score first.
	sortCandidates(cands)
	return cands
}

// filter applies hard constraints from the decision.
func filter(m reg.Model, d Decision) bool {
	if m.Availability != reg.Available {
		return false
	}
	if d.PreferredProvider != "" && m.Provider != d.PreferredProvider {
		return false
	}
	if d.RequiresCodeGeneration && !m.Capabilities.CodeGeneration {
		return false
	}
	// Budget cap: if we know pricing, enforce it for a nominal request.
	if d.MaxCost > 0 && m.Capabilities.SupportsAny() {
		est := m.Pricing.CostFor(1000, 500)
		if est > d.MaxCost {
			return false
		}
	}
	return true
}

// score ranks a candidate between 0 and 1 based on the decision priority.
func score(m reg.Model, d Decision) float64 {
	base := 0.5

	// Coding tasks strongly prefer big-pickle.
	if d.RequiresCodeGeneration && m.Provider == "opencode" && m.Model == BigPickleModel {
		base += 0.8
	}

	switch d.Priority {
	case "cost":
		// Cheaper is better; normalize using a small reference cost.
		ref := m.Pricing.CostFor(1000, 500) + 0.000001
		base += 1.0 / (1.0 + ref*100)
	case "quality":
		if m.Capabilities.CodeReasoning || m.Capabilities.CodeGeneration {
			base += 0.4
		}
	case "latency", "":
		base += 0.1
	}
	// Slight preference for code-capable / tool-capable models across priors.
	if m.Capabilities.Tools {
		base += 0.1
	}
	return base
}

func sortCandidates(cands []candidate) {
	for i := 1; i < len(cands); i++ {
		for j := i; j > 0 && cands[j].score > cands[j-1].score; j-- {
			cands[j], cands[j-1] = cands[j-1], cands[j]
		}
	}
}
