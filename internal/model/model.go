// Package model provides the central model registry and model metadata
// (Master Architecture §10).
package model

// Availability describes the operational state of a model.
type Availability string

const (
	Available   Availability = "available"
	Unavailable Availability = "unavailable"
	Deprecated  Availability = "deprecated"
)

// IDs uniquely identify a model by provider and model ID.
type IDs struct {
	Provider string
	Model    string
}

// Model is a registered model entry with its metadata.
type Model struct {
	IDs
	Description string
	// Capabilities of the model (vision, tools, streaming, etc).
	Capabilities Capabilities
	// Pricing per 1M tokens.
	Pricing Pricing
	// ContextWindow is the maximum input context size in tokens.
	ContextWindow int
	// MaxOutputTokens is the maximum completion length in tokens.
	MaxOutputTokens int
	// Availability of the model (routable or not).
	Availability Availability
}
