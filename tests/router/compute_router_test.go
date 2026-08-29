package routertest

import (
	"errors"
	"testing"

	cd "github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
	cr "github.com/zwbzd26dby-beep/Kitty-Go/internal/router/compute"
)

func compReg(devices ...cd.Device) *cd.Registry {
	r := cd.NewRegistry()
	for _, d := range devices {
		_ = r.Register(d)
	}
	return r
}

func gpuDevice(id string, vram, ram int, models ...string) cd.Device {
	return cd.Device{
		ID:              id,
		Resources:       cd.Resources{VRAMMB: vram, RAMMB: ram},
		SupportedModels: models,
		Health:          cd.HealthHealthy,
	}
}

func TestComputeRouterFiltersByVRAM(t *testing.T) {
	r := compReg(
		gpuDevice("small", 2000, 8000),
		gpuDevice("big", 16000, 64000, "llama3"),
	)
	rt := cr.NewRouter(r)
	d, err := rt.Select(cr.Request{NeedVRAM: 8000})
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if d.ID != "big" {
		t.Fatalf("expected big (only one with enough VRAM), got %q", d.ID)
	}
}

func TestComputeRouterFiltersByModel(t *testing.T) {
	r := compReg(
		gpuDevice("a", 8000, 32000, "llama3"),
		gpuDevice("b", 8000, 32000, "qwen"),
	)
	rt := cr.NewRouter(r)
	d, err := rt.Select(cr.Request{PreferredModel: "qwen"})
	if err != nil {
		t.Fatal(err)
	}
	if d.ID != "b" {
		t.Fatalf("expected b (serves qwen), got %q", d.ID)
	}
}

func TestComputeRouterExcludesUnhealthy(t *testing.T) {
	unhealthy := gpuDevice("down", 16000, 64000)
	unhealthy.Health = cd.HealthUnhealthy
	r := compReg(unhealthy)
	rt := cr.NewRouter(r)
	_, err := rt.Select(cr.Request{NeedVRAM: 1000})
	if !errors.Is(err, cr.ErrNoDeviceAvailable) {
		t.Fatalf("expected ErrNoDeviceAvailable, got %v", err)
	}
}

func TestComputeRouterNoDevices(t *testing.T) {
	rt := cr.NewRouter(compReg())
	if _, err := rt.Select(cr.Request{}); !errors.Is(err, cr.ErrNoDeviceAvailable) {
		t.Fatalf("expected ErrNoDeviceAvailable, got %v", err)
	}
}

func TestComputeRouterFallbackExcludes(t *testing.T) {
	r := compReg(gpuDevice("a", 8000, 16000), gpuDevice("b", 8000, 16000))
	rt := cr.NewRouter(r)
	d, err := rt.GetFallback(cr.Request{}, "a")
	if err != nil {
		t.Fatal(err)
	}
	if d.ID == "a" {
		t.Fatalf("expected fallback to exclude a, got %q", d.ID)
	}
}

type staticMetrics struct {
	loads   map[string]float64
	latency map[string]float64
}

func (m staticMetrics) Load(id string) float64      { return m.loads[id] }
func (m staticMetrics) LatencyMS(id string) float64 { return m.latency[id] }

func TestComputeRouterScoresByLoadAndLatency(t *testing.T) {
	r := compReg(gpuDevice("busy", 16000, 64000), gpuDevice("idle", 16000, 64000))
	rt := cr.NewRouter(r).WithMetrics(staticMetrics{
		loads:   map[string]float64{"busy": 0.9, "idle": 0.1},
		latency: map[string]float64{"busy": 200, "idle": 30},
	})
	d, err := rt.Select(cr.Request{})
	if err != nil {
		t.Fatal(err)
	}
	if d.ID != "idle" {
		t.Fatalf("expected idle (lower load/latency), got %q", d.ID)
	}
}

func TestComputeSchedulerRoundRobin(t *testing.T) {
	r := compReg(gpuDevice("a", 8000, 16000), gpuDevice("b", 8000, 16000), gpuDevice("c", 8000, 16000))
	s := cr.NewScheduler(r)
	seen := map[string]bool{}
	for i := 0; i < 3; i++ {
		d, err := s.Next(cr.Request{})
		if err != nil {
			t.Fatal(err)
		}
		seen[d.ID] = true
	}
	if len(seen) != 3 {
		t.Fatalf("expected round-robin across 3 devices, saw %d", len(seen))
	}
}
