package compute

import rc "github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"

// ResourceAware is a routing strategy that prioritises devices with the most
// available resources for the task.
type ResourceAware struct {
	router *Router
}

// NewResourceAware builds a ResourceAware selector over reg.
func NewResourceAware(reg DeviceLister) *ResourceAware {
	return &ResourceAware{router: NewRouter(reg)}
}

// Select returns the resource-best device for the request.
func (r *ResourceAware) Select(req Request) (rc.Device, error) {
	return r.router.Select(req)
}
