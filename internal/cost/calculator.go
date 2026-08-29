// Package cost implements cost calculation, tracking and alerting
// (Master Architecture §11).
package cost

// Usage is the token usage of a single completion.
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Pricing per token (or per 1M token denominator).
type Pricing struct {
	InputCostPerMillion  float64
	OutputCostPerMillion float64
}

// Record is an individual cost event for a request.
type Record struct {
	Provider string
	Model    string
	Usage    Usage
	Cost     float64
}

// Calculator computes the monetary cost of a completion from usage+pricing.
type Calculator struct{}

// Calculate returns cost for the given usage and pricing.
func (Calculator) Calculate(u Usage, p Pricing) float64 {
	return (float64(u.PromptTokens)/1e6)*p.InputCostPerMillion +
		(float64(u.CompletionTokens)/1e6)*p.OutputCostPerMillion
}
