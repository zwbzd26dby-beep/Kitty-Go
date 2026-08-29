package model

// Pricing holds per-1M-token costs for a model.
type Pricing struct {
	// InputCostPerMillion is the cost per 1M input (prompt) tokens.
	InputCostPerMillion float64
	// OutputCostPerMillion is the cost per 1M output (completion) tokens.
	OutputCostPerMillion float64
}

// CostFor computes the cost of a request given token usage.
func (p Pricing) CostFor(inTokens, outTokens int) float64 {
	return (float64(inTokens)/1e6)*p.InputCostPerMillion +
		(float64(outTokens)/1e6)*p.OutputCostPerMillion
}
