// Package compute implements the Compute Router (Master Architecture §14),
// which selects a device for a task from the device pool.
package compute

import (
	"errors"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
)

// ErrNoDeviceAvailable is returned when no device satisfies the request.
var ErrNoDeviceAvailable = errors.New("no device available for the given constraints")

// Request specifies what a task needs from a device (Master Arch §34).
type Request struct {
	NeedVRAM int
	NeedRAM  int
	// PreferredModel optionally restricts to devices serving the model.
	PreferredModel string
}

// ComputeRouter selects a device for a task (Master Architecture §14).
type ComputeRouter interface {
	Select(req Request) (compute.Device, error)
	GetFallback(req Request, exclude ...string) (compute.Device, error)
}

// Router implements ComputeRouter over a device registry.
type Router struct {
	registry DeviceLister
	metrics  MetricsProvider
}

// DeviceLister is the minimal registry surface the router needs.
type DeviceLister interface {
	List() []compute.Device
}

// MetricsProvider supplies per-device load/latency used for scoring
// (Master Architecture §34). It is optional; when nil, scoring relies on
// resources alone.
type MetricsProvider interface {
	Load(id string) float64
	LatencyMS(id string) float64
}

// NewRouter builds a ComputeRouter over reg.
func NewRouter(reg DeviceLister) *Router {
	return &Router{registry: reg}
}

// WithMetrics attaches a MetricsProvider to the router.
func (r *Router) WithMetrics(m MetricsProvider) *Router {
	r.metrics = m
	return r
}

// Select returns the highest-scoring device satisfying the request.
func (r *Router) Select(req Request) (compute.Device, error) {
	return r.GetFallback(req)
}

// GetFallback returns the best candidate excluding the given device IDs.
func (r *Router) GetFallback(req Request, exclude ...string) (compute.Device, error) {
	excluded := make(map[string]bool, len(exclude))
	for _, e := range exclude {
		excluded[e] = true
	}
	cands := r.candidates(req, excluded)
	if len(cands) == 0 {
		return compute.Device{}, ErrNoDeviceAvailable
	}
	best := cands[0]
	for _, c := range cands[1:] {
		if c.score > best.score {
			best = c
		}
	}
	return best.device, nil
}

// Chain returns the candidate devices in descending score order.
func (r *Router) Chain(req Request) []compute.Device {
	cands := r.candidates(req, nil)
	out := make([]compute.Device, 0, len(cands))
	for _, c := range cands {
		out = append(out, c.device)
	}
	return out
}

var _ ComputeRouter = (*Router)(nil)
