package model

// CostAware is a scoring strategy that weights price heavily. It is a thin
// wrapper around the router exposing a cost-first selection.
type CostAware struct {
	router *Router
}

// NewCostAware builds a CostAware selector over reg.
func NewCostAware(reg ModelRegistry) *CostAware {
	return &CostAware{router: NewRouter(reg)}
}

// Cheapest returns the cheapest admissible model for the decision.
func (c *CostAware) Cheapest(req RouteRequest) (Model, error) {
	d := req.Decision
	d.Priority = "cost"
	req.Decision = d
	return c.router.Select(req)
}
