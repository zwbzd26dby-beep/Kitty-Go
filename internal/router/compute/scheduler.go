package compute

import rc "github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"

// Scheduler assigns requests to devices, optionally distributing load via
// round-robin. It wraps a ComputeRouter for the base selection.
type Scheduler struct {
	router *Router
	next   int
}

// NewScheduler builds a Scheduler over reg.
func NewScheduler(reg DeviceLister) *Scheduler {
	return &Scheduler{router: NewRouter(reg)}
}

// Next selects a device, preferring round-robin among candidates when
// multiple are admissible.
func (s *Scheduler) Next(req Request) (rc.Device, error) {
	chain := s.router.Chain(req)
	if len(chain) == 0 {
		return rc.Device{}, ErrNoDeviceAvailable
	}
	if len(chain) == 1 {
		return chain[0], nil
	}
	idx := s.next % len(chain)
	s.next = (s.next + 1) % len(chain)
	return chain[idx], nil
}
