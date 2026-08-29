package computetest

import (
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
)

func testDevice(id string) compute.Device {
	return compute.Device{
		ID:              id,
		Name:            "node-" + id,
		Address:         "127.0.0.1:9000",
		SupportedModels: []string{"llama3"},
		Health:          compute.HealthUnknown,
	}
}

func TestRegistryRegisterGet(t *testing.T) {
	r := compute.NewRegistry()
	if err := r.Register(testDevice("a")); err != nil {
		t.Fatal(err)
	}
	d, err := r.Get("a")
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != "node-a" {
		t.Fatalf("unexpected name %q", d.Name)
	}
	if d.Health != compute.HealthUnknown {
		t.Fatalf("expected unknown health, got %q", d.Health)
	}
}

func TestRegistryRegisterRequiresID(t *testing.T) {
	r := compute.NewRegistry()
	if err := r.Register(compute.Device{}); err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestRegistryUnregister(t *testing.T) {
	r := compute.NewRegistry()
	r.Register(testDevice("a"))
	if err := r.Unregister("a"); err != nil {
		t.Fatal(err)
	}
	if len(r.List()) != 0 {
		t.Fatal("expected empty registry")
	}
	if err := r.Unregister("missing"); err == nil {
		t.Fatal("expected error unregistering missing")
	}
}

func TestRegistryListAndUpdateHealth(t *testing.T) {
	r := compute.NewRegistry()
	r.Register(testDevice("a"))
	r.Register(testDevice("b"))
	if got := len(r.List()); got != 2 {
		t.Fatalf("expected 2 devices, got %d", got)
	}
	if err := r.UpdateHealth("a", compute.HealthHealthy); err != nil {
		t.Fatal(err)
	}
	d, _ := r.Get("a")
	if d.Health != compute.HealthHealthy {
		t.Fatalf("expected healthy, got %q", d.Health)
	}
}

func TestFabricHealthyDevices(t *testing.T) {
	r := compute.NewRegistry()
	r.Register(testDevice("ok"))
	bad := testDevice("bad")
	bad.Health = compute.HealthUnhealthy
	r.Register(bad)
	f := compute.NewFabric(r, nil)
	healthy := f.HealthyDevices()
	if len(healthy) != 1 {
		t.Fatalf("expected 1 healthy device, got %d", len(healthy))
	}
	if healthy[0].ID != "ok" {
		t.Fatalf("expected ok device, got %q", healthy[0].ID)
	}
}
