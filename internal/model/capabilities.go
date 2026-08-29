package model

// Capabilities describes what a model can do. Missing capabilities default to
// false, which helps filters exclude models for unsupported tasks.
type Capabilities struct {
	CodeGeneration bool
	CodeReasoning  bool
	Vision         bool
	Tools          bool
	Streaming      bool
	FunctionCall   bool
}

// SupportsAny returns true if the model advertises at least one capability,
// used as a coarse "has metadata" check.
func (c Capabilities) SupportsAny() bool {
	return c.CodeGeneration || c.CodeReasoning || c.Vision || c.Tools || c.Streaming || c.FunctionCall
}
