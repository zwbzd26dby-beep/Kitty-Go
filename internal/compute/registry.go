package compute

import (
	"fmt"
	"sync"
)

// DeviceRegistry stores and indexes compute devices (Master Architecture §16).
type DeviceRegistry interface {
	Register(d Device) error
	Unregister(id string) error
	Get(id string) (Device, error)
	List() []Device
	UpdateHealth(id string, h HealthState) error
}

// Registry is a thread-safe in-memory DeviceRegistry.
type Registry struct {
	mu      sync.RWMutex
	devices map[string]Device
}

// NewRegistry creates an empty DeviceRegistry.
func NewRegistry() *Registry {
	return &Registry{devices: make(map[string]Device)}
}

// Register adds or replaces a device.
func (r *Registry) Register(d Device) error {
	if d.ID == "" {
		return fmt.Errorf("device id is required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.Health == "" {
		d.Health = HealthUnknown
	}
	r.devices[d.ID] = d
	return nil
}

// Unregister removes a device.
func (r *Registry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.devices[id]; !ok {
		return fmt.Errorf("device %q not found", id)
	}
	delete(r.devices, id)
	return nil
}

// Get returns a device by id.
func (r *Registry) Get(id string) (Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.devices[id]
	if !ok {
		return Device{}, fmt.Errorf("device %q not found", id)
	}
	return d, nil
}

// List returns all devices.
func (r *Registry) List() []Device {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Device, 0, len(r.devices))
	for _, d := range r.devices {
		out = append(out, d)
	}
	return out
}

// UpdateHealth updates a device's health state.
func (r *Registry) UpdateHealth(id string, h HealthState) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.devices[id]
	if !ok {
		return fmt.Errorf("device %q not found", id)
	}
	d.Health = h
	r.devices[id] = d
	return nil
}

var _ DeviceRegistry = (*Registry)(nil)
