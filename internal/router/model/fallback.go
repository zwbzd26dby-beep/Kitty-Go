package model

// Chain returns the available candidates in descending score order, giving
// callers the full ordered list to walk for progressive fallback.
func (r *Router) Chain(req RouteRequest) []Model {
	cands := r.candidates(req, nil)
	out := make([]Model, 0, len(cands))
	for _, c := range cands {
		out = append(out, c.model)
	}
	return out
}
