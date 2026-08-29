package compute

import (
	"net"
	"time"
)

// Pinger reports whether a device is reachable/healthy.
type Pinger interface {
	Ping(d Device) error
}

// TCPPinger checks health by opening a TCP connection to the device address.
type TCPPinger struct {
	Timeout time.Duration
}

// Ping dials the device address.
func (p TCPPinger) Ping(d Device) error {
	timeout := p.Timeout
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	conn, err := net.DialTimeout("tcp", d.Address, timeout)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}

// HealthChecker periodically pings all devices and updates their state.
type HealthChecker struct {
	pinger Pinger
}

// NewHealthChecker creates a HealthChecker using pinger.
func NewHealthChecker(pinger Pinger) *HealthChecker {
	return &HealthChecker{pinger: pinger}
}

// Start begins periodic health checks. Each element of tick triggers a full
// pass; the loop exits when tick is closed.
func (h *HealthChecker) Start(reg DeviceRegistry, tick <-chan struct{}) error {
	for range tick {
		for _, d := range reg.List() {
			state := HealthHealthy
			if h.pinger != nil {
				if err := h.pinger.Ping(d); err != nil {
					state = HealthUnhealthy
				}
			}
			_ = reg.UpdateHealth(d.ID, state)
		}
	}
	return nil
}

// RunOnce performs a single health-check pass synchronously (for tests).
func (h *HealthChecker) RunOnce(reg DeviceRegistry) error {
	tick := make(chan struct{}, 1)
	tick <- struct{}{}
	close(tick)
	return h.Start(reg, tick)
}
