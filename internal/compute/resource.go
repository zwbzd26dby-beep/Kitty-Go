package compute

// resourceAdequate reports whether a device has enough resources for a task
// requiring the given VRAM/RAM.
func resourceAdequate(d Device, needVRAM, needRAM int) bool {
	if needVRAM > 0 && d.Resources.VRAMMB < needVRAM {
		return false
	}
	if needRAM > 0 && d.Resources.RAMMB < needRAM {
		return false
	}
	return true
}
