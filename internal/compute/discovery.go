package compute

// staticDiscover registers the statically-configured devices. In this phase
// devices are seeded directly through the registry (e.g. from config); the
// static pass is intentionally additive and non-destructive.
func staticDiscover(reg DeviceRegistry) error {
	// No hardcoded static devices by default; callers seed via config.
	return nil
}

// LANDiscoverer implements device discovery over the local network using a
// broadcast UDP probe. It is a lightweight mechanism independent of mDNS.
type LANDiscoverer struct {
	// ProbePort is the port workers listen on.
	ProbePort int
	// InterfaceAddr is the local broadcast address (e.g. 255.255.255.255).
	InterfaceAddr string
}

// Discover emits a broadcast probe and returns discovered device addresses.
// Real probing is left to Phase 9/10 wiring; this scaffolds the contract.
func (d *LANDiscoverer) Discover() ([]string, error) {
	return nil, nil
}
