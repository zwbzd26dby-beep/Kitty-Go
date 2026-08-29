package computetest

import (
	"errors"
	"testing"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
)

type fakePinger struct {
	fail map[string]bool
}

func (f fakePinger) Ping(d compute.Device) error {
	if f.fail[d.ID] {
		return errors.New("unreachable")
	}
	return nil
}

func TestHealthCheckerMarksUnhealthy(t *testing.T) {
	r := compute.NewRegistry()
	r.Register(testDevice("ok"))
	bad := testDevice("bad")
	r.Register(bad)

	hc := compute.NewHealthChecker(fakePinger{fail: map[string]bool{"bad": true}})
	if err := hc.RunOnce(r); err != nil {
		t.Fatal(err)
	}
	d1, _ := r.Get("ok")
	d2, _ := r.Get("bad")
	if d1.Health != compute.HealthHealthy {
		t.Fatalf("expected healthy for ok, got %q", d1.Health)
	}
	if d2.Health != compute.HealthUnhealthy {
		t.Fatalf("expected unhealthy for bad, got %q", d2.Health)
	}
}

func TestHealthCheckerPeriodic(t *testing.T) {
	r := compute.NewRegistry()
	r.Register(testDevice("a"))
	hc := compute.NewHealthChecker(fakePinger{})

	tick := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = hc.Start(r, tick)
	}()
	tick <- struct{}{}
	tick <- struct{}{}
	close(tick)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("health loop did not stop")
	}
	d, _ := r.Get("a")
	if d.Health != compute.HealthHealthy {
		t.Fatalf("expected healthy after periodic check, got %q", d.Health)
	}
}

func TestFabricStartsHealthMonitor(t *testing.T) {
	r := compute.NewRegistry()
	r.Register(testDevice("a"))
	hc := compute.NewHealthChecker(fakePinger{})
	f := compute.NewFabric(r, hc)

	tick := make(chan struct{})
	go func() {
		tick <- struct{}{}
		close(tick)
	}()
	if err := f.StartHealthMonitor(tick); err != nil {
		t.Fatal(err)
	}
	d, _ := r.Get("a")
	if d.Health != compute.HealthHealthy {
		t.Fatalf("expected healthy after fabric monitor, got %q", d.Health)
	}
}
