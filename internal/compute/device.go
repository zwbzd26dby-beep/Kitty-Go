// Package compute implements the Compute Fabric and Device Registry
// (Master Architecture §15-16), later extended with distributed execution.
package compute

// Resources describes the compute capacity of a device.
type Resources struct {
	// VRAMMB is available GPU memory in megabytes.
	VRAMMB int
	// RAMMB is available system memory in megabytes.
	RAMMB int
	// CPUCount is the number of logical CPUs.
	CPUCount int
	// GPUName is the GPU model, if any.
	GPUName string
}

// HealthState is the operational status of a device.
type HealthState string

const (
	HealthUnknown   HealthState = "unknown"
	HealthHealthy   HealthState = "healthy"
	HealthUnhealthy HealthState = "unhealthy"
	HealthOffline   HealthState = "offline"
)

// Device is a single compute node in the fabric.
type Device struct {
	// ID is a stable unique identifier.
	ID string
	// Name is a human-friendly name.
	Name string
	// Address is the reachable endpoint (host:port).
	Address string
	// Resources of the device.
	Resources Resources
	// SupportedModels the device can serve.
	SupportedModels []string
	// AuthToken is used to authenticate execution (never logged).
	AuthToken string
	// Health of the device.
	Health HealthState
}
