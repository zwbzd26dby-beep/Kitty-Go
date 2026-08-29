package compute

// ComputeFabric is the high-level facade over the device pool
// (Master Architecture §15). It combines registration, discovery, health and
// selection.
type ComputeFabric interface {
	Registry() DeviceRegistry
	// Discover launches static + dynamic discovery.
	Discover() error
	// StartHealthMonitor begins periodic health checks.
	StartHealthMonitor(pollInterval <-chan struct{}) error
	// HealthyDevices returns all currently healthy devices.
	HealthyDevices() []Device
}

// Fabric is the default ComputeFabric implementation.
type Fabric struct {
	registry DeviceRegistry
	checker  *HealthChecker
}

// NewFabric builds a Fabric over reg with the given health checker.
func NewFabric(reg DeviceRegistry, checker *HealthChecker) *Fabric {
	return &Fabric{registry: reg, checker: checker}
}

// Registry returns the underlying DeviceRegistry.
func (f *Fabric) Registry() DeviceRegistry { return f.registry }

// Discover performs static registration (already seeded) and dynamic
// discovery. See discovery.go for the LAN discovery.
func (f *Fabric) Discover() error {
	return staticDiscover(f.registry)
}

// StartHealthMonitor runs periodic health checks until the poll channel is
// closed. It runs in its own goroutine; errors do not stop the loop.
func (f *Fabric) StartHealthMonitor(tick <-chan struct{}) error {
	if f.checker == nil {
		return nil
	}
	f.checker.Start(f.registry, tick)
	return nil
}

// HealthyDevices returns only devices with HealthHealthy state.
func (f *Fabric) HealthyDevices() []Device {
	var out []Device
	for _, d := range f.registry.List() {
		if d.Health == HealthHealthy || d.Health == HealthUnknown {
			out = append(out, d)
		}
	}
	return out
}

var _ ComputeFabric = (*Fabric)(nil)
